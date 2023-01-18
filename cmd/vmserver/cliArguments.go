package main

import (
	"github.com/multiversx/wasm-vm/vmserver"
	"github.com/urfave/cli"
)

type cliArguments struct {
	// Common arguments
	ServerAddress string
	Database      string
	World         string
	Outcome       string
	// For contract-related actions
	Impersonated    string
	ContractAddress string
	Action          string
	Function        string
	Arguments       cli.StringSlice
	Code            string
	CodePath        string
	CodeMetadata    string
	Value           string
	GasLimit        uint64
	GasPrice        uint64
	// For blockchain-related action
	AccountAddress string
	AccountBalance string
	AccountNonce   uint64
}

func (args *cliArguments) toDeployRequest() vmserver.DeployRequest {
	request := &vmserver.DeployRequest{}
	args.populateDeployRequest(request)

	return *request
}

func (args *cliArguments) populateDeployRequest(request *vmserver.DeployRequest) {
	args.populateContractRequestBase(&request.ContractRequestBase)

	request.CodeHex = args.Code
	request.CodePath = args.CodePath
	request.CodeMetadata = args.CodeMetadata
	request.ArgumentsHex = args.Arguments
}

func (args *cliArguments) populateContractRequestBase(request *vmserver.ContractRequestBase) {
	args.populateRequestBase(&request.RequestBase)

	request.ImpersonatedHex = args.Impersonated
	request.Value = args.Value
	request.GasLimit = args.GasLimit
	request.GasPrice = args.GasPrice
}

func (args *cliArguments) populateRequestBase(request *vmserver.RequestBase) {
	request.DatabasePath = args.Database
	request.World = args.World
	request.Outcome = args.Outcome
}

func (args *cliArguments) toUpgradeRequest() vmserver.UpgradeRequest {
	request := &vmserver.UpgradeRequest{}
	args.populateDeployRequest(&request.DeployRequest)

	request.ContractAddressHex = args.ContractAddress
	return *request
}

func (args *cliArguments) toRunRequest() vmserver.RunRequest {
	request := &vmserver.RunRequest{}
	args.populateRunRequest(request)

	return *request
}

func (args *cliArguments) populateRunRequest(request *vmserver.RunRequest) {
	args.populateContractRequestBase(&request.ContractRequestBase)

	request.ContractAddressHex = args.ContractAddress
	request.Function = args.Function
	request.ArgumentsHex = args.Arguments
}

func (args *cliArguments) toQueryRequest() vmserver.QueryRequest {
	request := &vmserver.QueryRequest{}
	args.populateRunRequest(&request.RunRequest)

	return *request
}

func (args *cliArguments) toCreateAccountRequest() vmserver.CreateAccountRequest {
	request := &vmserver.CreateAccountRequest{}
	args.populateRequestBase(&request.RequestBase)

	request.AddressHex = args.AccountAddress
	request.Balance = args.AccountBalance
	request.Nonce = args.AccountNonce
	return *request
}
