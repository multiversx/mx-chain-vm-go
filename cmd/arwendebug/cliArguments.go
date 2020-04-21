package main

import "github.com/ElrondNetwork/arwen-wasm-vm/arwendebug"

type cliArguments struct {
	ServerAddress   string
	DatabasePath    string
	Session         string
	Impersonator    string
	ContractAddress string
	Action          string
	Function        string
	Arguments       []string
	Code            string
	CodePath        string
	CodeMetadata    string
}

func (args *cliArguments) toDeployRequest() arwendebug.DeployRequest {
	request := &arwendebug.DeployRequest{}
	args.populateDeployRequest(request)
	return *request
}

func (args *cliArguments) populateDeployRequest(request *arwendebug.DeployRequest) {
	args.populateRequestBase(&request.RequestBase)

	request.Code = args.Code
	request.CodeMetadata = args.CodeMetadata
	request.Arguments = args.Arguments
}

func (args *cliArguments) populateRequestBase(request *arwendebug.RequestBase) {
	request.DatabasePath = args.DatabasePath
	request.Session = args.Session
	request.Impersonator = args.Impersonator
}

func (args *cliArguments) toUpgradeRequest() arwendebug.UpgradeRequest {
	request := &arwendebug.UpgradeRequest{}
	request.ContractAddress = args.ContractAddress
	args.populateDeployRequest(&request.DeployRequest)
	return *request
}
