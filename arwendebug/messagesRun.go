package arwendebug

// RunRequest is a CLI / REST request message
type RunRequest struct {
	ContractRequestBase
	ContractAddressAsHex string
	Function             string
	ArgumentsAsHex       []string
	ArgumentsAsBytes     [][]byte
}

func (request *RunRequest) digest() error {
	err := request.ContractRequestBase.digest()
	if err != nil {
		return err
	}

	request.ArgumentsAsBytes, err = decodeArguments(request.ArgumentsAsHex)
	if err != nil {
		return err
	}

	return nil
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
