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

type EIFunction struct {
	OriginalName string
	PublicName   string
	Arguments    []*EIFunctionArg
	Result       *EIFunctionResult
}

type EIFileMetadata struct {
	// TODO: add a name and imports function name
	Functions []*EIFunction
}

type EIMetadata struct {
	FileMap      map[string]*EIFileMetadata
	AllFunctions []*EIFunction
}
