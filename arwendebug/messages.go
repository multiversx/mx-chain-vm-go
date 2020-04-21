package arwendebug

import (
	"encoding/json"

	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

// RequestBase -
type RequestBase struct {
	DatabasePath string
	Session      string
	Impersonator string
}

// ResponseBase -
type ResponseBase struct {
	Input  vmcommon.VMInput
	Output vmcommon.VMOutput
	Error  error
}

// DebugString -
func (response *ResponseBase) DebugString() string {
	data, err := json.MarshalIndent(response, "", "\t")
	if err != nil {
		return "{}"
	}

	return string(data)
}

// DeployRequest -
type DeployRequest struct {
	RequestBase
	Code         string
	CodeMetadata string
	Arguments    []string
}

// DeployResponse -
type DeployResponse struct {
	ResponseBase
}

// UpgradeRequest -
type UpgradeRequest struct {
	DeployRequest
	ContractAddress string
}

// RunRequest -
type RunRequest struct {
	RequestBase
	ContractAddress string
	Function        string
	Arguments       []string
}

// QueryRequest -
type QueryRequest struct {
	RunRequest
}
