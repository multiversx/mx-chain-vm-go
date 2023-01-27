package vmserver

// RunRequest is a CLI / REST request message
type RunRequest struct {
	ContractRequestBase
	ContractAddressHex string
	ContractAddress    []byte
	Function           string
	ArgumentsHex       []string
	Arguments          [][]byte
}

func (request *RunRequest) digest() error {
	err := request.ContractRequestBase.digest()
	if err != nil {
		return err
	}

	request.Arguments, err = decodeArguments(request.ArgumentsHex)
	if err != nil {
		return err
	}

	request.ContractAddress, err = fromHex(request.ContractAddressHex)
	if err != nil {
		return err
	}

	return nil
}

// RunResponse is a CLI / REST response message
type RunResponse struct {
	ContractResponseBase
}
