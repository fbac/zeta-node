// @generated by protoc-gen-es v1.3.0 with parameter "target=dts"
// @generated from file pkg/proofs/proofs.proto (package proofs, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Message, proto3 } from "@bufbuild/protobuf";
import type { Proof as Proof$1 } from "./ethereum/ethereum_pb.js";
import type { Proof as Proof$2 } from "./bitcoin/bitcoin_pb.js";

/**
 * @generated from message proofs.BlockHeader
 */
export declare class BlockHeader extends Message<BlockHeader> {
  /**
   * @generated from field: int64 height = 1;
   */
  height: bigint;

  /**
   * @generated from field: bytes hash = 2;
   */
  hash: Uint8Array;

  /**
   * @generated from field: bytes parent_hash = 3;
   */
  parentHash: Uint8Array;

  /**
   * @generated from field: int64 chain_id = 4;
   */
  chainId: bigint;

  /**
   * chain specific header
   *
   * @generated from field: proofs.HeaderData header = 5;
   */
  header?: HeaderData;

  constructor(data?: PartialMessage<BlockHeader>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "proofs.BlockHeader";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): BlockHeader;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): BlockHeader;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): BlockHeader;

  static equals(a: BlockHeader | PlainMessage<BlockHeader> | undefined, b: BlockHeader | PlainMessage<BlockHeader> | undefined): boolean;
}

/**
 * @generated from message proofs.HeaderData
 */
export declare class HeaderData extends Message<HeaderData> {
  /**
   * @generated from oneof proofs.HeaderData.data
   */
  data: {
    /**
     * binary encoded headers; RLP for ethereum
     *
     * @generated from field: bytes ethereum_header = 1;
     */
    value: Uint8Array;
    case: "ethereumHeader";
  } | {
    /**
     * 80-byte little-endian encoded binary data
     *
     * @generated from field: bytes bitcoin_header = 2;
     */
    value: Uint8Array;
    case: "bitcoinHeader";
  } | { case: undefined; value?: undefined };

  constructor(data?: PartialMessage<HeaderData>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "proofs.HeaderData";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): HeaderData;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): HeaderData;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): HeaderData;

  static equals(a: HeaderData | PlainMessage<HeaderData> | undefined, b: HeaderData | PlainMessage<HeaderData> | undefined): boolean;
}

/**
 * @generated from message proofs.Proof
 */
export declare class Proof extends Message<Proof> {
  /**
   * @generated from oneof proofs.Proof.proof
   */
  proof: {
    /**
     * @generated from field: ethereum.Proof ethereum_proof = 1;
     */
    value: Proof$1;
    case: "ethereumProof";
  } | {
    /**
     * @generated from field: bitcoin.Proof bitcoin_proof = 2;
     */
    value: Proof$2;
    case: "bitcoinProof";
  } | { case: undefined; value?: undefined };

  constructor(data?: PartialMessage<Proof>);

  static readonly runtime: typeof proto3;
  static readonly typeName = "proofs.Proof";
  static readonly fields: FieldList;

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): Proof;

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): Proof;

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): Proof;

  static equals(a: Proof | PlainMessage<Proof> | undefined, b: Proof | PlainMessage<Proof> | undefined): boolean;
}
