package arwendebug

import (
	"io/ioutil"

	"github.com/multiversx/mx-chain-vm-common-go"
)

// DeployRequest is a CLI / REST request message
type DeployRequest struct {
	ContractRequestBase
	CodeHex           string
	Code              []byte
	CodePath          string
	CodeMetadata      string
	CodeMetadataBytes []byte
	ArgumentsHex      []string
	Arguments         [][]byte
}

func (request *DeployRequest) digest() error {
	err := request.ContractRequestBase.digest()
	if err != nil {
		return err
	}

	if len(request.CodeHex) > 0 {
		request.Code, err = fromHex(request.CodeHex)
		if err != nil {
			return NewRequestErrorMessageInner("invalid contract code", err)
		}
	}

	if len(request.CodePath) > 0 {
		request.Code, err = ioutil.ReadFile(request.CodePath)
		if err != nil {
			return err
		}
	}

	if len(request.Code) == 0 {
		return NewRequestError("invalid contract code")
	}

	request.CodeMetadataBytes = (&vmcommon.CodeMetadata{Upgradeable: true}).ToBytes()
	if len(request.CodeMetadata) > 0 {
		request.CodeMetadataBytes, err = fromHex(request.CodeMetadata)
		if err != nil {
			return err
		}
	}

	request.Arguments, err = decodeArguments(request.ArgumentsHex)
	if err != nil {
		return err
	}

	return nil
}

// DeployResponse is a CLI / REST response message
type DeployResponse struct {
	ContractResponseBase
	ContractAddress    []byte
	ContractAddressHex string
}
