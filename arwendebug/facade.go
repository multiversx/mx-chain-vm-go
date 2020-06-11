package arwendebug

import (
	"encoding/json"
	"fmt"

	logger "github.com/ElrondNetwork/elrond-go-logger"
)

var log = logger.GetOrCreate("arwendebug")

// DebugFacade is the debug facade
type DebugFacade struct {
}

// NewDebugFacade creates a new debug facade
func NewDebugFacade() {
}

// DeploySmartContract deploys a smart contract
func (facade *DebugFacade) DeploySmartContract(request DeployRequest) (*DeployResponse, error) {
	log.Debug("DebugFacade.DeploySmartContract()")

	database := facade.loadDatabase(request.DatabasePath)
	world, err := database.loadWorld(request.World)
	if err != nil {
		return nil, err
	}

	response, err := world.deploySmartContract(request)
	if err != nil {
		return nil, err
	}

	err = database.storeWorld(world)
	if err != nil {
		return nil, err
	}

	err = database.storeOutcome(request.Outcome, response)
	if err != nil {
		return nil, err
	}

	dumpOutcome(&response)
	return response, err
}

func (facade *DebugFacade) loadDatabase(rootPath string) *database {
	database := newDatabase(rootPath)
	return database
}

// UpgradeSmartContract upgrades a smart contract
func (facade *DebugFacade) UpgradeSmartContract(request UpgradeRequest) (*UpgradeResponse, error) {
	log.Debug("DebugFacade.UpgradeSmartContract()")

	database := facade.loadDatabase(request.DatabasePath)
	world, err := database.loadWorld(request.World)
	if err != nil {
		return nil, err
	}

	response, err := world.upgradeSmartContract(request)
	if err != nil {
		return nil, err
	}

	err = database.storeWorld(world)
	if err != nil {
		return nil, err
	}

	err = database.storeOutcome(request.Outcome, response)
	if err != nil {
		return nil, err
	}

	dumpOutcome(&response)
	return response, err
}

// RunSmartContract executes a smart contract function
func (facade *DebugFacade) RunSmartContract(request RunRequest) (*RunResponse, error) {
	log.Debug("DebugFacade.RunSmartContract()")

	database := facade.loadDatabase(request.DatabasePath)
	world, err := database.loadWorld(request.World)
	if err != nil {
		return nil, err
	}

	response, err := world.runSmartContract(request)
	if err != nil {
		return nil, err
	}

	err = database.storeWorld(world)
	if err != nil {
		return nil, err
	}

	err = database.storeOutcome(request.Outcome, response)
	if err != nil {
		return nil, err
	}

	dumpOutcome(&response)
	return response, err
}

// QuerySmartContract queries a pure function of the smart contract
func (facade *DebugFacade) QuerySmartContract(request QueryRequest) (*QueryResponse, error) {
	log.Debug("DebugFacade.QuerySmartContracts()")

	database := facade.loadDatabase(request.DatabasePath)
	world, err := database.loadWorld(request.World)
	if err != nil {
		return nil, err
	}

	response, err := world.querySmartContract(request)
	if err != nil {
		return nil, err
	}

	err = database.storeOutcome(request.Outcome, response)
	if err != nil {
		return nil, err
	}

	dumpOutcome(&response)
	return response, err
}

// CreateAccount creates a test account
func (facade *DebugFacade) CreateAccount(request CreateAccountRequest) (*CreateAccountResponse, error) {
	log.Debug("DebugFacade.CreateAccount()")

	database := facade.loadDatabase(request.DatabasePath)
	world, err := database.loadWorld(request.World)
	if err != nil {
		return nil, err
	}

	response, err := world.qreateAccount(request)
	if err != nil {
		return nil, err
	}

	err = database.storeWorld(world)
	if err != nil {
		return nil, err
	}

	err = database.storeOutcome(request.Outcome, response)
	if err != nil {
		return nil, err
	}

	dumpOutcome(&response)
	return response, err
}

func dumpOutcome(outcome interface{}) {
	data, err := json.Marshal(outcome)
	if err != nil {
		fmt.Println("{}")
	}

	fmt.Println(string(data))
}
