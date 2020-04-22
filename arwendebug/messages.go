package arwendebug

import (
	"math/big"

	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

// RequestBase -
type RequestBase struct {
	DatabasePath string
	World        string
}

// ResponseBase -
type ResponseBase struct {
	Error error
}

// ContractRequestBase -
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

// ContractResponseBase -
type ContractResponseBase struct {
	ResponseBase
	Input  *vmcommon.VMInput
	Output *vmcommon.VMOutput
}

// DeployRequest -
type DeployRequest struct {
	ContractRequestBase
	Code         string
	CodeMetadata string
	Arguments    []string
}

// DeployResponse -
type DeployResponse struct {
	ContractResponseBase
}

// UpgradeRequest -
type UpgradeRequest struct {
	DeployRequest
	ContractAddress string
}

// UpgradeResponse -
type UpgradeResponse struct {
	ContractResponseBase
}

// RunRequest -
type RunRequest struct {
	ContractRequestBase
	ContractAddress string
	Function        string
	Arguments       []string
}

// RunResponse -
type RunResponse struct {
	ContractResponseBase
}

// QueryRequest -
type QueryRequest struct {
	RunRequest
}

// QueryResponse -
type QueryResponse struct {
	ContractResponseBase
}

// CreateAccountRequest -
type CreateAccountRequest struct {
	RequestBase
	Address string
	Balance string
	Nonce   uint64
}

func (request *CreateAccountRequest) parseBalance() (*big.Int, error) {
	balance, ok := big.NewInt(0).SetString(request.Balance, 10)
	if !ok {
		return nil, NewRequestError("invalid balance")
	}

	return balance, nil
}

// CreateAccountResponse -
type CreateAccountResponse struct {
}
