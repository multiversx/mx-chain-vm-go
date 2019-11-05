package config

type BaseOperationCost struct {
	StorePerByte    uint64
	DataCopyPerByte uint64
}

type ElrondAPICost struct {
	GetOwner           uint64
	GetExternalBalance uint64
	GetBlockHash       uint64
	TransferValue      uint64
	GetArgument        uint64
	GetFunction        uint64
	GetNumArguments    uint64
	StorageStore       uint64
	StorageLoad        uint64
	GetCaller          uint64
	GetCallValue       uint64
	Log                uint64
	Finish             uint64
	SignalError        uint64
	GetBlockTimeStamp  uint64
	GetGasLeft         uint64
	Int64GetArgument   uint64
	Int64StorageStore  uint64
	Int64StorageLoad   uint64
	Int64Finish        uint64
}

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
	BigIntNew                uint64
	BigIntByteLength         uint64
	BigIntGetBytes           uint64
	BigIntSetBytes           uint64
	BigIntIsInt64            uint64
	BigIntGetInt64           uint64
	BigIntSetInt64           uint64
	BigIntAdd                uint64
	BigIntSub                uint64
	BigIntMul                uint64
	BigIntCmp                uint64
	BigIntFinish             uint64
	BigIntStorageLoad        uint64
	BigIntStorageStore       uint64
	BigIntGetArgument        uint64
	BigIntGetCallValue       uint64
	BigIntGetExternalBalance uint64
}

type GasCost struct {
	BaseOperationCost BaseOperationCost
	BigIntAPICost     BigIntAPICost
	EthAPICost        EthAPICost
	ElrondAPICost     ElrondAPICost
}
