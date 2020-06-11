package arwendebug

import (
	"encoding/hex"
	"io/ioutil"

	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

// DeployRequest is a CLI / REST request message
type DeployRequest struct {
	ContractRequestBase
	CodeAsHex           string
	CodeAsBytes         []byte
	CodePath            string
	CodeMetadata        string
	CodeMetadataAsBytes []byte
	ArgumentsAsHex      []string
	ArgumentsAsBytes    [][]byte
}

func (request *DeployRequest) digest() error {
	err := request.ContractRequestBase.digest()
	if err != nil {
		return err
	}

	if len(request.CodeAsHex) > 0 {
		request.CodeAsBytes, err = hex.DecodeString(request.CodeAsHex)
		if err != nil {
			return NewRequestErrorMessageInner("invalid contract code", err)
		}
	}

	if len(request.CodePath) > 0 {
		request.CodeAsBytes, err = ioutil.ReadFile(request.CodePath)
		if err != nil {
			return err
		}
	}

	if len(request.CodeAsBytes) == 0 {
		return NewRequestError("invalid contract code")
	}

	request.CodeMetadataAsBytes = (&vmcommon.CodeMetadata{Upgradeable: true}).ToBytes()
	if len(request.CodeMetadata) > 0 {
		request.CodeMetadataAsBytes, err = hex.DecodeString(request.CodeMetadata)
		if err != nil {
			return err
		}
	}

	request.ArgumentsAsBytes, err = decodeArguments(request.ArgumentsAsHex)
	if err != nil {
		return err
	}

	return nil
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
