package local

import (
	"fmt"
	"runtime"
	"time"

	"github.com/fatih/color"

	"github.com/zeta-chain/zetacore/e2e/config"
	"github.com/zeta-chain/zetacore/e2e/e2etests"
	"github.com/zeta-chain/zetacore/e2e/runner"
)

// zevmMPTestRoutine runs ZEVM message passing related e2e tests
func zevmMPTestRoutine(
	conf config.Config,
	deployerRunner *runner.E2ERunner,
	verbose bool,
	testNames ...string,
) func() error {
	return func() (err error) {
		// return an error on panic
		// TODO: remove and instead return errors in the tests
		// https://github.com/zeta-chain/node/issues/1500
		defer func() {
			if r := recover(); r != nil {
				// print stack trace
				stack := make([]byte, 4096)
				n := runtime.Stack(stack, false)
				err = fmt.Errorf("zevm mp panic: %v, stack trace %s", r, stack[:n])
			}
		}()

		// initialize runner for zevm mp test
		zevmMPRunner, err := initTestRunner(
			"zevm_mp",
			conf,
			deployerRunner,
			UserZEVMMPTestAddress,
			UserZEVMMPTestPrivateKey,
			runner.NewLogger(verbose, color.FgHiRed, "zevm_mp"),
		)
		if err != nil {
			return err
		}

		zevmMPRunner.Logger.Print("🏃 starting ZEVM Message Passing tests")
		startTime := time.Now()

		// funding the account
		txZetaSend := deployerRunner.SendZetaOnEvm(UserZEVMMPTestAddress, 1000)
		zevmMPRunner.WaitForTxReceiptOnEvm(txZetaSend)

		// depositing the necessary tokens on ZetaChain
		txZetaDeposit := zevmMPRunner.DepositZeta()
		txEtherDeposit := zevmMPRunner.DepositEther(false)
		zevmMPRunner.WaitForMinedCCTX(txZetaDeposit)
		zevmMPRunner.WaitForMinedCCTX(txEtherDeposit)

		// run zevm message passing test
		testsToRun, err := zevmMPRunner.GetE2ETestsToRunByName(
			e2etests.AllE2ETests,
			testNames...,
		)
		if err != nil {
			return fmt.Errorf("zevm message passing tests failed: %v", err)
		}

		if err := zevmMPRunner.RunE2ETests(testsToRun); err != nil {
			return fmt.Errorf("zevm message passing tests failed: %v", err)
		}

		zevmMPRunner.Logger.Print("🍾 ZEVM message passing tests completed in %s", time.Since(startTime).String())

		return err
	}
}