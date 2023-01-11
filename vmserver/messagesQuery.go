package vmserver

// QueryRequest is a CLI / REST request message
type QueryRequest struct {
	RunRequest
}

// QueryResponse is a CLI / REST response message
type QueryResponse struct {
	ContractResponseBase
}
