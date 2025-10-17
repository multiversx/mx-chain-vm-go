package evmhooks

func (context *EVMHooksImpl) GasLeft() uint64 {
	return context.GetMeteringContext().GasLeft()
}

func (context *EVMHooksImpl) UseGas(opCode string, gas uint64) bool {
	err := context.GetMeteringContext().UseGasBoundedAndAddTracedGas(opCode, gas)
	return err == nil
}

func (context *EVMHooksImpl) BlockGasLimit() uint64 {
	return context.GetMeteringContext().BlockGasLimit()
}
