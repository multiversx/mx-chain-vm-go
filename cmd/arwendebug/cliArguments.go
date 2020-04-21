package main

import "github.com/ElrondNetwork/arwen-wasm-vm/arwendebug"

type cliArguments struct {
	// Common arguments
	ServerAddress string
	Database      string
	Session       string
	// For contract-related actions
	Impersonated    string
	ContractAddress string
	Action          string
	Function        string
	Arguments       []string
	Code            string
	CodePath        string
	CodeMetadata    string
	// For blockchain-related action
	AccountAddress string
	AccountBalance string
	AccountNonce   uint64
}

func (args *cliArguments) toDeployRequest() arwendebug.DeployRequest {
	request := &arwendebug.DeployRequest{}
	args.populateDeployRequest(request)
	return *request
}

func (args *cliArguments) populateDeployRequest(request *arwendebug.DeployRequest) {
	args.populateContractRequestBase(&request.ContractRequestBase)

	request.Code = args.Code
	request.CodeMetadata = args.CodeMetadata
	request.Arguments = args.Arguments
}

func (args *cliArguments) populateContractRequestBase(request *arwendebug.ContractRequestBase) {
	args.populateRequestBase(&request.RequestBase)
	request.Impersonated = args.Impersonated
}

func (args *cliArguments) populateRequestBase(request *arwendebug.RequestBase) {
	request.DatabasePath = args.Database
	request.Session = args.Session
}

func (args *cliArguments) toUpgradeRequest() arwendebug.UpgradeRequest {
	request := &arwendebug.UpgradeRequest{}
	request.ContractAddress = args.ContractAddress
	args.populateDeployRequest(&request.DeployRequest)
	return *request
}

func (args *cliArguments) toCreateAccountRequest() arwendebug.CreateAccountRequest {
	request := &arwendebug.CreateAccountRequest{}
	args.populateRequestBase(&request.RequestBase)

	request.Balance = args.AccountBalance
	return *request
}
