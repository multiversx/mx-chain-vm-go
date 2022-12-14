package elrondapigenerate

// EIFunctionArg models an executor callback method arg.
type EIFunctionArg struct {
	Name string
	Type string
}

// EIFunctionResult models the executor callback method result.
type EIFunctionResult struct {
	Type string
}

// EIFunction holds data about one function in the VM EI.
type EIFunction struct {
	Name      string
	Arguments []*EIFunctionArg
	Result    *EIFunctionResult
}

// EIGroup groups EI functions into bundles.
// They can end up in separate interfaces or files, if desired.
type EIGroup struct {
	SourcePath string
	Name       string
	Functions  []*EIFunction
}

// EIMetadata holds all data about EI functions in the VM.
type EIMetadata struct {
	Groups       []*EIGroup
	AllFunctions []*EIFunction
}
