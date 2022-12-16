package contexts

// GasTraceMap is the map that holds gas traces for all API functions, depending on SCAddress and functionName
type gasTraceMap map[string]map[string][]uint64

// gasTracer holds information about the gas trace
type gasTracer struct {
	gasTrace           gasTraceMap
	functionNameTraced string
	scAddress          string
}

// NewEnabledGasTracer creates a new gasTracer
func NewEnabledGasTracer() *gasTracer {
	return &gasTracer{
		gasTrace: make(map[string]map[string][]uint64)}
}

// NewDisabledGasTracer creates a new disabledGasTracer
func NewDisabledGasTracer() *disabledGasTracer {
	return &disabledGasTracer{}
}

// BeginTrace prepares the gasTracer to trace the correct scAddress and function
func (gt *gasTracer) BeginTrace(scAddress string, functionName string) {
	gt.setCurrentFunctionTraced(functionName)
	gt.setCurrentScAddressTraced(scAddress)
	gt.createGasTraceIfNil(scAddress, functionName)
	gt.newTrace(scAddress, functionName)
}

// AddToCurrentTrace ads the usedGas passed at the current trace value
func (gt *gasTracer) AddToCurrentTrace(usedGas uint64) {
	gt.createGasTraceIfNil(gt.scAddress, gt.functionNameTraced)
	funcTrace := gt.gasTrace[gt.scAddress][gt.functionNameTraced]
	length := len(funcTrace)
	if length > 0 {
		gt.gasTrace[gt.scAddress][gt.functionNameTraced][length-1] += usedGas
	} else {
		gt.gasTrace[gt.scAddress][gt.functionNameTraced] = append(funcTrace, usedGas)
	}
}

// AddTracedGas directly ads usedGas in the gasTrace map
func (gt *gasTracer) AddTracedGas(scAddress string, functionName string, usedGas uint64) {
	gt.createGasTraceIfNil(scAddress, functionName)
	gt.gasTrace[scAddress][functionName] = append(gt.gasTrace[scAddress][functionName], usedGas)
}

func (gt *gasTracer) setCurrentFunctionTraced(functionName string) {
	gt.functionNameTraced = functionName
}

func (gt *gasTracer) setCurrentScAddressTraced(scAddress string) {
	gt.scAddress = scAddress
}

func (gt *gasTracer) newTrace(scAddress string, functionName string) {
	gt.gasTrace[scAddress][functionName] = append(gt.gasTrace[scAddress][functionName], 0)
}

func (gt *gasTracer) createGasTraceIfNil(scAddress string, functionName string) {
	gt.createSCAdressGasTracingIfNil(scAddress)
	gt.createFunctionNameGasTracingIfNil(scAddress, functionName)
}

func (gt *gasTracer) createSCAdressGasTracingIfNil(scAddress string) {
	if gt.gasTrace[scAddress] == nil {
		gt.gasTrace[scAddress] = make(map[string][]uint64)
	}
}

func (gt *gasTracer) createFunctionNameGasTracingIfNil(scAddress string, functionName string) {
	if gt.gasTrace[scAddress][functionName] == nil {
		gt.gasTrace[scAddress][functionName] = make([]uint64, 0)
	}
}

// GetGasTrace returns the gasTrace map
func (gt *gasTracer) GetGasTrace() map[string]map[string][]uint64 {
	return gt.gasTrace
}

// IsInterfaceNil returns true if there is no value under the interface
func (gt *gasTracer) IsInterfaceNil() bool {
	return gt == nil
}

type disabledGasTracer struct {
}

// BeginTrace does nothing
func (dgt *disabledGasTracer) BeginTrace(_ string, _ string) {
}

// AddToCurrentTrace does nothing
func (dgt *disabledGasTracer) AddToCurrentTrace(_ uint64) {
}

// AddTracedGas does nothing
func (dgt *disabledGasTracer) AddTracedGas(_ string, _ string, _ uint64) {
}

// GetGasTrace returns nil
func (dgt *disabledGasTracer) GetGasTrace() map[string]map[string][]uint64 {
	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (dgt *disabledGasTracer) IsInterfaceNil() bool {
	return dgt == nil
}
