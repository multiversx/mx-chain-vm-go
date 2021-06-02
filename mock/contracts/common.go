package contracts

// DirectCallGasTestConfig is configuration for direct call tests
type DirectCallGasTestConfig struct {
	GasUsedByParent      uint64
	GasUsedByChild       uint64
	GasProvidedToChild   uint64
	GasProvided          uint64
	ParentBalance        int64
	ChildBalance         int64
	ESDTTokensToTransfer uint64
}

// AsyncCallBaseTestConfig is base configuration for async call tests
type AsyncCallBaseTestConfig struct {
	GasProvided       uint64
	GasUsedByParent   uint64
	GasUsedByChild    uint64
	GasUsedByCallback uint64
	GasLockCost       uint64

	TransferFromParentToChild int64

	ParentBalance int64
	ChildBalance  int64
}

// AsyncCallTestConfig is configuration for async call tests
type AsyncCallTestConfig struct {
	AsyncCallBaseTestConfig
	TransferToThirdParty         int64
	TransferToVault              int64
	ESDTTokensToTransfer         uint64
	CallbackESDTTokensToTransfer uint64
}

// AsyncBuiltInCallTestConfig is configuration for async call tests of builtin functions
type AsyncBuiltInCallTestConfig struct {
	AsyncCallBaseTestConfig
	TransferFromChildToParent int64
}

// AsyncCallRecursiveTestConfig is configuration for recursive async call tests
type AsyncCallRecursiveTestConfig struct {
	AsyncCallBaseTestConfig
	RecursiveChildCalls int
}

// AsyncCallMultiChildTestConfig is configuration for recursivemultiple children async call tests
type AsyncCallMultiChildTestConfig struct {
	AsyncCallBaseTestConfig
	ChildCalls int
}

// GasTestConfig interface for gas tests configs
type GasTestConfig interface {
	GetGasUsedByChild() uint64
}

// GetGasUsedByChild - getter for GasUsedByChild
func (config AsyncCallTestConfig) GetGasUsedByChild() uint64 {
	return config.GasUsedByChild
}

// GetGasUsedByChild - getter for GasUsedByChild
func (config DirectCallGasTestConfig) GetGasUsedByChild() uint64 {
	return config.GasUsedByChild
}
