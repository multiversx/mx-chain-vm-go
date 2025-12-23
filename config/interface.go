package config

import (
	"github.com/multiversx/mx-chain-vm-go/executor"
)

type GasSchedule interface {
	GetBaseOperationCost() *BaseOperationCost
	GetBigIntAPICost() *BigIntAPICost
	GetBigFloatAPICost() *BigFloatAPICost
	GetBaseOpsAPICost() *BaseOpsAPICost
	GetManagedBufferAPICost() *ManagedBufferAPICost
	GetManagedMapAPICost() *ManagedMapAPICost
	GetCryptoAPICost() *CryptoAPICost
	GetWASMOpcodeCost() *executor.WASMOpcodeCost
	GetDynamicStorageLoad() *DynamicStorageLoadCostCoefficients
}

type GasScheduleFactory interface {
	CreateGasSchedule(gasMap GasScheduleMap) (GasSchedule, error)
	IsInterfaceNil() bool
}
