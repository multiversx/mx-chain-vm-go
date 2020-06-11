package arwendebug

import (
	"encoding/hex"
	"math/big"

	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

// RequestBase is a CLI / REST request message
type RequestBase struct {
	DatabasePath string
	World        string
	Outcome      string
}

// ResponseBase is a CLI / REST response message
type ResponseBase struct {
	Error error
}

// ContractRequestBase is a CLI / REST request message
type ContractRequestBase struct {
	RequestBase
	Impersonated string
	Value        string
	GasPrice     uint64
	GasLimit     uint64
}

func (request *ContractRequestBase) getValue() *big.Int {
	value := big.NewInt(0)
	_, _ = value.SetString(request.Value, 10)
	return value
}

func (request *ContractRequestBase) getGasPrice() uint64 {
	if request.GasPrice == 0 {
		return DefaultGasPrice
	}

	return request.GasPrice
}

func (request *ContractRequestBase) getGasLimit() uint64 {
	if request.GasLimit == 0 {
		return DefaultGasLimit
	}

	return request.GasLimit
}

func (request *ContractRequestBase) getImpersonated() ([]byte, error) {
	if request.Impersonated == "" {
		return nil, NewRequestError("empty impersonated address")
	}

	impersonatedAsHex := request.Impersonated
	impersonatedAsBytes, err := hex.DecodeString(impersonatedAsHex)
	if err != nil {
		return nil, NewRequestErrorMessageInner("invalid impersonated address", err)
	}

	return impersonatedAsBytes, nil
}

// ContractResponseBase is a CLI / REST response message
type ContractResponseBase struct {
	ResponseBase
	Input  *vmcommon.VMInput
	Output *vmcommon.VMOutput
}
