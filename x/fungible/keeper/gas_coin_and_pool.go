package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/zetacore/common"
	contracts "github.com/zeta-chain/zetacore/contracts/zevm"
	"github.com/zeta-chain/zetacore/x/fungible/types"
	zetaObserverTypes "github.com/zeta-chain/zetacore/x/observer/types"
	"math/big"
)

// setup gas ZRC20, and ZETA/gas pool for a chain
// add 0.1gas/0.1wzeta to the pool
// FIXME: use chainid instead of chain name; add cointype and use proper gas limit based on cointype/chain
func (k Keeper) setupChainGasCoinAndPool(ctx sdk.Context, c string, gasAssetName string, symbol string, decimals uint8) (ethcommon.Address, error) {
	name := fmt.Sprintf("%s-%s", gasAssetName, c)
	chainName := common.ParseChainName(c)
	chain := k.zetaobserverKeeper.GetParams(ctx).GetChainFromChainName(chainName)
	if chain == nil {
		return ethcommon.Address{}, zetaObserverTypes.ErrSupportedChains
	}

	transferGasLimit := big.NewInt(21_000)
	if chain.IsBitcoinChain() {
		transferGasLimit = big.NewInt(100) // 100B for a typical tx
	}

	zrc20Addr, err := k.DeployZRC20Contract(ctx, name, symbol, decimals, chain.ChainName.String(), common.CoinType_Gas, "", transferGasLimit)
	if err != nil {
		return ethcommon.Address{}, sdkerrors.Wrapf(err, "failed to DeployZRC20Contract")
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(name, zrc20Addr.String()),
		),
	)
	err = k.SetGasCoin(ctx, big.NewInt(chain.ChainId), zrc20Addr)
	if err != nil {
		return ethcommon.Address{}, err
	}
	amount := big.NewInt(10)
	amount.Exp(amount, big.NewInt(int64(decimals-1)), nil)
	amountAZeta := big.NewInt(1e17)

	_, err = k.DepositZRC20(ctx, zrc20Addr, types.ModuleAddressEVM, amount)
	if err != nil {
		return ethcommon.Address{}, err
	}
	err = k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin("azeta", sdk.NewIntFromBigInt(amountAZeta))))
	if err != nil {
		return ethcommon.Address{}, err
	}
	systemContractAddress, err := k.GetSystemContractAddress(ctx)
	if err != nil || systemContractAddress == (ethcommon.Address{}) {
		return ethcommon.Address{}, sdkerrors.Wrapf(types.ErrContractNotFound, "system contract address invalid: %s", systemContractAddress)
	}
	systemABI, err := contracts.SystemContractMetaData.GetAbi()
	if err != nil {
		return ethcommon.Address{}, sdkerrors.Wrapf(err, "failed to get system contract abi")
	}
	_, err = k.CallEVM(ctx, *systemABI, types.ModuleAddressEVM, systemContractAddress, BigIntZero, nil, true, "setGasZetaPool", big.NewInt(chain.ChainId), zrc20Addr)
	if err != nil {
		return ethcommon.Address{}, sdkerrors.Wrapf(err, "failed to CallEVM method setGasZetaPool(%d, %s)", chain.ChainId, zrc20Addr.String())
	}

	// setup uniswap v2 pools gas/zeta
	routerAddress, err := k.GetUniswapV2Router02Address(ctx)
	if err != nil {
		return ethcommon.Address{}, sdkerrors.Wrapf(err, "failed to GetUniswapV2Router02Address")
	}
	routerABI, err := contracts.UniswapV2Router02MetaData.GetAbi()
	if err != nil {
		return ethcommon.Address{}, sdkerrors.Wrapf(err, "failed to get uniswap router abi")
	}
	zrc4ABI, err := contracts.ZRC20MetaData.GetAbi()
	if err != nil {
		return ethcommon.Address{}, sdkerrors.Wrapf(err, "failed to GetAbi zrc20")
	}
	_, err = k.CallEVM(ctx, *zrc4ABI, types.ModuleAddressEVM, zrc20Addr, BigIntZero, nil, true, "approve", routerAddress, amount)
	if err != nil {
		return ethcommon.Address{}, sdkerrors.Wrapf(err, "failed to CallEVM method approve(%s, %d)", routerAddress.String(), amount)
	}
	//function addLiquidityETH(
	//	address token,
	//	uint amountTokenDesired,
	//	uint amountTokenMin,
	//	uint amountETHMin,
	//	address to,
	//	uint deadline
	//) external payable returns (uint amountToken, uint amountETH, uint liquidity);
	res, err := k.CallEVM(ctx, *routerABI, types.ModuleAddressEVM, routerAddress, amount, big.NewInt(20_000_000), true,
		"addLiquidityETH", zrc20Addr, amount, BigIntZero, BigIntZero, types.ModuleAddressEVM, amountAZeta)
	if err != nil {
		return ethcommon.Address{}, sdkerrors.Wrapf(err, "failed to CallEVM method addLiquidityETH(%s, %s)", zrc20Addr.String(), amountAZeta.String())
	}
	AmountToken := new(*big.Int)
	AmountETH := new(*big.Int)
	Liquidity := new(*big.Int)
	err = routerABI.UnpackIntoInterface(&[]interface{}{AmountToken, AmountETH, Liquidity}, "addLiquidityETH", res.Ret)
	if err != nil {
		return ethcommon.Address{}, sdkerrors.Wrapf(err, "failed to UnpackIntoInterface addLiquidityETH")

	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute("function", "addLiquidityETH"),
			sdk.NewAttribute("amountToken", (*AmountToken).String()),
			sdk.NewAttribute("amountETH", (*AmountETH).String()),
			sdk.NewAttribute("liquidity", (*Liquidity).String()),
		),
	)
	return zrc20Addr, nil
}