package arwendebug

import (
	"encoding/hex"
	"io/ioutil"

	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

// DeployRequest is a CLI / REST request message
type DeployRequest struct {
	ContractRequestBase
	Code         string
	CodePath     string
	CodeMetadata string
	Arguments    []string
}

func (request *DeployRequest) getCode() ([]byte, error) {
	if len(request.Code) > 0 {
		codeAsHex := request.Code
		codeAsBytes, err := hex.DecodeString(codeAsHex)
		if err != nil {
			return nil, NewRequestErrorMessageInner("invalid contract code", err)
		}

		return codeAsBytes, nil
	}

	if len(request.CodePath) > 0 {
		codeAsBytes, err := ioutil.ReadFile(request.CodePath)
		if err != nil {
			return nil, err
		}

		return codeAsBytes, nil
	}

	return nil, NewRequestError("invalid contract code")
}

func (request *DeployRequest) getCodeMetadata() ([]byte, error) {
	if len(request.CodeMetadata) > 0 {
		metadataAsHex := request.CodeMetadata
		metadataAsBytes, err := hex.DecodeString(metadataAsHex)
		if err != nil {
			return nil, err
		}

		return metadataAsBytes, nil
	}

	defaultMetadata := vmcommon.CodeMetadata{Upgradeable: true}
	return defaultMetadata.ToBytes(), nil
}

func (request *DeployRequest) getArguments() ([][]byte, error) {
	return decodeArguments(request.Arguments)
}

// DeployResponse is a CLI / REST response message
type DeployResponse struct {
	ContractResponseBase
	ContractAddress string
}

// UpgradeRequest is a CLI / REST request message
type UpgradeRequest struct {
	DeployRequest
	ContractAddress string
}

// UpgradeResponse is a CLI / REST response message
type UpgradeResponse struct {
	ContractResponseBase
}
