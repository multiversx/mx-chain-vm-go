package arwendebug

// RunRequest is a CLI / REST request message
type RunRequest struct {
	ContractRequestBase
	ContractAddress string
	Function        string
	Arguments       []string
}

func (request *RunRequest) getArguments() ([][]byte, error) {
	return decodeArguments(request.Arguments)
}

// RunResponse is a CLI / REST response message
type RunResponse struct {
	ContractResponseBase
}

// QueryRequest is a CLI / REST request message
type QueryRequest struct {
	RunRequest
}

// QueryResponse is a CLI / REST response message
type QueryResponse struct {
	ContractResponseBase
}
