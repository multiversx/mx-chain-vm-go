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

func (request *RequestBase) digest() error {
	if request.DatabasePath == "" {
		request.DatabasePath = "./db"
	}

	if request.World == "" {
		request.World = "default"
	}

	return nil
}

// ResponseBase is a CLI / REST response message
type ResponseBase struct {
	Error error
}

// ContractRequestBase is a CLI / REST request message
type ContractRequestBase struct {
	RequestBase
	ImpersonatedAsHex   string
	ImpersonatedAsBytes []byte
	Value               string
	ValueAsBigInt       *big.Int
	GasPrice            uint64
	GasLimit            uint64
}

func (request *ContractRequestBase) digest() error {
	var err error
	var ok bool

	err = request.RequestBase.digest()
	if err != nil {
		return err
	}

	if request.ImpersonatedAsHex == "" {
		return NewRequestError("empty impersonated address")
	}

	request.ImpersonatedAsBytes, err = hex.DecodeString(request.ImpersonatedAsHex)
	if err != nil {
		return NewRequestErrorMessageInner("invalid impersonated address", err)
	}

	if request.GasPrice == 0 {
		request.GasPrice = DefaultGasPrice
	}

	if request.GasLimit == 0 {
		request.GasLimit = DefaultGasLimit
	}

	// todo move to common
	request.ValueAsBigInt = big.NewInt(0)
	if len(request.Value) > 0 {
		_, ok = request.ValueAsBigInt.SetString(request.Value, 10)
		if !ok {
			return NewRequestError("invalid value (erd)")
		}
	}

	return nil
}

// ContractResponseBase is a CLI / REST response message
type ContractResponseBase struct {
	ResponseBase
	Input  *vmcommon.VMInput
	Output *vmcommon.VMOutput
}
