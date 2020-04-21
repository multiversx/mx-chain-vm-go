package arwendebug

// RequestBase -
type RequestBase struct {
	DatabasePath string
	Session      string
	Impersonator string
}

// DeployRequest -
type DeployRequest struct {
	RequestBase
	Code         string
	CodeMetadata string
	Arguments    []string
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
