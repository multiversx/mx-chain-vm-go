package contexts

// GasTraceMap is the map that holds gas traces for all API functions, depending on SCAddress and functionName
type gasTraceMap map[string]map[string][]uint64

// GasTracer holds information about the gas trace
type GasTracer struct {
	traceGasEnabled    bool
	gasTrace           gasTraceMap
	functionNameTraced string
	scAddress          string
}

// NewGasTracer creates a new GasTracer
func NewGasTracer() *GasTracer {
	return &GasTracer{
		traceGasEnabled: false,
		gasTrace:        make(map[string]map[string][]uint64)}
}

// BeginTrace prepares the gasTracer to trace the correct scAddress and function
func (gt *GasTracer) BeginTrace(scAddress string, functionName string) {
	gt.setCurrentFunctionTraced(functionName)
	gt.setCurrentScAddressTraced(scAddress)
	gt.createGasTraceIfNil(scAddress, functionName)
	gt.newTrace(scAddress, functionName)
}

// AddToCurrentTrace ads the usedGas passed at the current trace value
func (gt *GasTracer) AddToCurrentTrace(usedGas uint64) {
	length := len(gt.gasTrace[gt.scAddress][gt.functionNameTraced])
	gt.gasTrace[gt.scAddress][gt.functionNameTraced][length-1] += usedGas
}

// AddTracedGas directly ads usedGas in the gasTrace map
func (gt *GasTracer) AddTracedGas(scAddress string, functionName string, usedGas uint64) {
	gt.createGasTraceIfNil(scAddress, functionName)
	gt.gasTrace[scAddress][functionName] = append(gt.gasTrace[scAddress][functionName], usedGas)

}

// SetTraceGasEnabled enables or disables gas tracing depending on the passed value
func (gt *GasTracer) SetTraceGasEnabled(setValue bool) {
	gt.traceGasEnabled = setValue
}

// IsEnabled returns true if gas tracing is enabled, false if it's disabled
func (gt *GasTracer) IsEnabled() bool {
	return gt.traceGasEnabled
}

// IsInterfaceNil returns true if there is no value under the interface
func (gt *GasTracer) IsInterfaceNil() bool {
	return gt == nil
}

func (gt *GasTracer) setCurrentFunctionTraced(functionName string) {
	gt.functionNameTraced = functionName
}

func (gt *GasTracer) setCurrentScAddressTraced(scAddress string) {
	gt.scAddress = scAddress
}

func (gt *GasTracer) newTrace(scAddress string, functionName string) {
	gt.gasTrace[scAddress][functionName] = append(gt.gasTrace[scAddress][functionName], 0)
}

func (gt *GasTracer) createGasTraceIfNil(scAddress string, functionName string) {
	gt.createSCAdressGasTracingIfNil(scAddress)
	gt.createFunctionNameGasTracingIfNil(scAddress, functionName)
}

func (gt *GasTracer) createSCAdressGasTracingIfNil(scAddress string) {
	if gt.gasTrace[scAddress] == nil {
		gt.gasTrace[scAddress] = make(map[string][]uint64)
	}
}

func (gt *GasTracer) createFunctionNameGasTracingIfNil(scAddress string, functionName string) {
	if gt.gasTrace[scAddress][functionName] == nil {
		gt.gasTrace[scAddress][functionName] = make([]uint64, 0)
	}
}

func (gt *GasTracer) GetGasTrace() map[string]map[string][]uint64 {
	return gt.gasTrace
}
