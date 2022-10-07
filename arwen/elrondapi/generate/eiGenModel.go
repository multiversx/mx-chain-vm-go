package elrondapigenerate

type EIFunctionValue int

const (
	EIFunctionValueInt32 EIFunctionValue = iota
	EIFunctionValueInt64
)

type EIFunctionArg struct {
	Name string
	Type string
}

type EIFunctionResult struct {
	Type string
}

// EIFunction holds data about one function in the VM EI.
type EIFunction struct {
	OriginalName string
	PublicName   string
	Arguments    []*EIFunctionArg
	Result       *EIFunctionResult
}

// EIGroup groups EI functions into bundles.
// They end up in separate interfaces.
type EIGroup struct {
	SourcePath string
	// TODO: add a name and imports function name
	Functions []*EIFunction
}

// EIMetadata holds all data about EI functions in the VM.
type EIMetadata struct {
	Groups       []*EIGroup
	AllFunctions []*EIFunction
}
