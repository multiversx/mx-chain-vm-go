package arwendebug

import (
	"encoding/json"
	"fmt"

	logger "github.com/ElrondNetwork/elrond-go-logger"
)

var log = logger.GetOrCreate("arwendebug")

// DebugFacade -
type DebugFacade struct {
}

// NewDebugFacade -
func NewDebugFacade() {
}

// DeploySmartContract -
func (facade *DebugFacade) DeploySmartContract(request DeployRequest) (*DeployResponse, error) {
	log.Debug("DebugFacade.DeploySmartContract()")

	database := facade.loadDatabase(request.DatabasePath)
	world, err := database.loadWorld(request.World)
	if err != nil {
		return nil, err
	}

	response, err := world.DeploySmartContract(request)
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
	// TODO: use factory
	database := NewDatabase(rootPath)
	return database
}

// UpgradeSmartContract -
func (facade *DebugFacade) UpgradeSmartContract(request UpgradeRequest) (*UpgradeResponse, error) {
	log.Debug("DebugFacade.UpgradeSmartContract()")

	database := facade.loadDatabase(request.DatabasePath)
	world, err := database.loadWorld(request.World)
	if err != nil {
		return nil, err
	}

	response, err := world.UpgradeSmartContract()
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

// RunSmartContract -
func (facade *DebugFacade) RunSmartContract(request RunRequest) (*RunResponse, error) {
	log.Debug("DebugFacade.RunSmartContract()")

	database := facade.loadDatabase(request.DatabasePath)
	world, err := database.loadWorld(request.World)
	if err != nil {
		return nil, err
	}

	response, err := world.RunSmartContract(request)
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

// QuerySmartContract -
func (facade *DebugFacade) QuerySmartContract(request QueryRequest) (*QueryResponse, error) {
	log.Debug("DebugFacade.QuerySmartContracts()")

	database := facade.loadDatabase(request.DatabasePath)
	world, err := database.loadWorld(request.World)
	if err != nil {
		return nil, err
	}

	response, err := world.QuerySmartContract(request)
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

// CreateAccount -
func (facade *DebugFacade) CreateAccount(request CreateAccountRequest) (*CreateAccountResponse, error) {
	log.Debug("DebugFacade.CreateAccount()")

	database := facade.loadDatabase(request.DatabasePath)
	world, err := database.loadWorld(request.World)
	if err != nil {
		return nil, err
	}

	response, err := world.CreateAccount(request)
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
