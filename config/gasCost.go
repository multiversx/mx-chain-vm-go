package config

type GasCost struct {
	BaseOperationCost    BaseOperationCost
	BigIntAPICost        BigIntAPICost
	BigFloatAPICost      BigFloatAPICost
	EthAPICost           EthAPICost
	ElrondAPICost        ElrondAPICost
	ManagedBufferAPICost ManagedBufferAPICost
	CryptoAPICost        CryptoAPICost
	WASMOpcodeCost       WASMOpcodeCost
}

type BaseOperationCost struct {
	StorePerByte      uint64
	ReleasePerByte    uint64
	DataCopyPerByte   uint64
	PersistPerByte    uint64
	CompilePerByte    uint64
	AoTPreparePerByte uint64
	GetCode           uint64
}

type ElrondAPICost struct {
	GetSCAddress            uint64
	GetOwnerAddress         uint64
	IsSmartContract         uint64
	GetShardOfAddress       uint64
	GetExternalBalance      uint64
	GetBlockHash            uint64
	GetOriginalTxHash       uint64
	GetCurrentTxHash        uint64
	GetPrevTxHash           uint64
	TransferValue           uint64
	GetArgument             uint64
	GetFunction             uint64
	GetNumArguments         uint64
	StorageStore            uint64
	StorageLoad             uint64
	CachedStorageLoad       uint64
	GetCaller               uint64
	GetCallValue            uint64
	Log                     uint64
	Finish                  uint64
	SignalError             uint64
	GetBlockTimeStamp       uint64
	GetGasLeft              uint64
	Int64GetArgument        uint64
	Int64StorageStore       uint64
	Int64StorageLoad        uint64
	Int64Finish             uint64
	GetStateRootHash        uint64
	GetBlockNonce           uint64
	GetBlockEpoch           uint64
	GetBlockRound           uint64
	GetBlockRandomSeed      uint64
	ExecuteOnSameContext    uint64
	ExecuteOnDestContext    uint64
	DelegateExecution       uint64
	ExecuteReadOnly         uint64
	AsyncCallStep           uint64
	AsyncCallbackGasLock    uint64
	CreateAsyncCall         uint64
	SetAsyncCallback        uint64
	SetAsyncGroupCallback   uint64
	SetAsyncContextCallback uint64
	GetCallbackClosure      uint64
	CreateContract          uint64
	GetReturnData           uint64
	GetNumReturnData        uint64
	GetReturnDataSize       uint64
	CleanReturnData         uint64
	DeleteFromReturnData    uint64
}

// TODO remove this struct
type EthAPICost struct {
	UseGas              uint64
	GetAddress          uint64
	GetExternalBalance  uint64
	GetBlockHash        uint64
	Call                uint64
	CallDataCopy        uint64
	GetCallDataSize     uint64
	CallCode            uint64
	CallDelegate        uint64
	CallStatic          uint64
	StorageStore        uint64
	StorageLoad         uint64
	GetCaller           uint64
	GetCallValue        uint64
	CodeCopy            uint64
	GetCodeSize         uint64
	GetBlockCoinbase    uint64
	Create              uint64
	GetBlockDifficulty  uint64
	ExternalCodeCopy    uint64
	GetExternalCodeSize uint64
	GetGasLeft          uint64
	GetBlockGasLimit    uint64
	GetTxGasPrice       uint64
	Log                 uint64
	GetBlockNumber      uint64
	GetTxOrigin         uint64
	Finish              uint64
	Revert              uint64
	GetReturnDataSize   uint64
	ReturnDataCopy      uint64
	SelfDestruct        uint64
	GetBlockTimeStamp   uint64
}

type BigIntAPICost struct {
	BigIntNew                  uint64
	BigIntUnsignedByteLength   uint64
	BigIntSignedByteLength     uint64
	BigIntGetUnsignedBytes     uint64
	BigIntGetSignedBytes       uint64
	BigIntSetUnsignedBytes     uint64
	BigIntSetSignedBytes       uint64
	BigIntIsInt64              uint64
	BigIntGetInt64             uint64
	BigIntSetInt64             uint64
	BigIntAdd                  uint64
	BigIntSub                  uint64
	BigIntMul                  uint64
	BigIntSqrt                 uint64
	BigIntPow                  uint64
	BigIntLog                  uint64
	BigIntTDiv                 uint64
	BigIntTMod                 uint64
	BigIntEDiv                 uint64
	BigIntEMod                 uint64
	BigIntAbs                  uint64
	BigIntNeg                  uint64
	BigIntSign                 uint64
	BigIntCmp                  uint64
	BigIntNot                  uint64
	BigIntAnd                  uint64
	BigIntOr                   uint64
	BigIntXor                  uint64
	BigIntShr                  uint64
	BigIntShl                  uint64
	BigIntFinishUnsigned       uint64
	BigIntFinishSigned         uint64
	BigIntStorageLoadUnsigned  uint64
	BigIntStorageStoreUnsigned uint64
	BigIntGetUnsignedArgument  uint64
	BigIntGetSignedArgument    uint64
	BigIntGetCallValue         uint64
	BigIntGetExternalBalance   uint64
	CopyPerByteForTooBig       uint64
}

type BigFloatAPICost struct {
	BigFloatNewFromParts uint64
	BigFloatAdd          uint64
	BigFloatSub          uint64
	BigFloatMul          uint64
	BigFloatDiv          uint64
	BigFloatTruncate     uint64
	BigFloatNeg          uint64
	BigFloatClone        uint64
	BigFloatCmp          uint64
	BigFloatAbs          uint64
	BigFloatSqrt         uint64
	BigFloatPow          uint64
	BigFloatFloor        uint64
	BigFloatCeil         uint64
	BigFloatIsInt        uint64
	BigFloatSetBigInt    uint64
	BigFloatSetInt64     uint64
	BigFloatGetConst     uint64
}

type CryptoAPICost struct {
	SHA256                 uint64
	Keccak256              uint64
	Ripemd160              uint64
	VerifyBLS              uint64
	VerifyEd25519          uint64
	VerifySecp256k1        uint64
	EllipticCurveNew       uint64
	AddECC                 uint64
	DoubleECC              uint64
	IsOnCurveECC           uint64
	ScalarMultECC          uint64
	MarshalECC             uint64
	MarshalCompressedECC   uint64
	UnmarshalECC           uint64
	UnmarshalCompressedECC uint64
	GenerateKeyECC         uint64
	EncodeDERSig           uint64
}

type ManagedBufferAPICost struct {
	MBufferNew                uint64
	MBufferNewFromBytes       uint64
	MBufferGetLength          uint64
	MBufferGetBytes           uint64
	MBufferGetByteSlice       uint64
	MBufferCopyByteSlice      uint64
	MBufferSetBytes           uint64
	MBufferAppend             uint64
	MBufferAppendBytes        uint64
	MBufferToBigIntUnsigned   uint64
	MBufferToBigIntSigned     uint64
	MBufferFromBigIntUnsigned uint64
	MBufferFromBigIntSigned   uint64
	MBufferToBigFloat         uint64
	MBufferFromBigFloat       uint64
	MBufferStorageStore       uint64
	MBufferStorageLoad        uint64
	MBufferGetArgument        uint64
	MBufferFinish             uint64
	MBufferSetRandom          uint64
}
