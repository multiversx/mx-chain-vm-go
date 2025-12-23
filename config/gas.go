package config

import (
	"github.com/multiversx/mx-chain-vm-go/executor"
)

var _ GasSchedule = (*GasCost)(nil)

func (g *GasCost) GetBaseOperationCost() *BaseOperationCost {
	return &g.BaseOperationCost
}

func (g *GasCost) GetBigIntAPICost() *BigIntAPICost {
	return &g.BigIntAPICost
}

func (g *GasCost) GetBigFloatAPICost() *BigFloatAPICost {
	return &g.BigFloatAPICost
}

func (g *GasCost) GetBaseOpsAPICost() *BaseOpsAPICost {
	return &g.BaseOpsAPICost
}

func (g *GasCost) GetManagedBufferAPICost() *ManagedBufferAPICost {
	return &g.ManagedBufferAPICost
}

func (g *GasCost) GetManagedMapAPICost() *ManagedMapAPICost {
	return &g.ManagedMapAPICost
}

func (g *GasCost) GetCryptoAPICost() *CryptoAPICost {
	return &g.CryptoAPICost
}

func (g *GasCost) GetWASMOpcodeCost() *executor.WASMOpcodeCost {
	return g.WASMOpcodeCost
}

func (g *GasCost) GetDynamicStorageLoad() *DynamicStorageLoadCostCoefficients {
	return &g.DynamicStorageLoad
}
