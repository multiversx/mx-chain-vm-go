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

	world, err := facade.loadWorld(request.DatabasePath, request.World)
	if err != nil {
		return nil, err
	}

	response, err := world.DeploySmartContract(request)
	if err != nil {
		return nil, err
	}

	dumpResponse(&response)
	return response, err
}

func (facade *DebugFacade) loadWorld(databaseRootPath string, worldID string) (*world, error) {
	database := facade.loadDatabase(databaseRootPath)
	return database.loadWorld(worldID)
}

func (facade *DebugFacade) loadDatabase(rootPath string) *database {
	// TODO: use factory
	database := NewDatabase(rootPath)
	return database
}

func (facade *DebugFacade) saveWorld() {

}

// UpgradeSmartContract -
func (facade *DebugFacade) UpgradeSmartContract(request UpgradeRequest) (*UpgradeResponse, error) {
	log.Debug("DebugFacade.UpgradeSmartContract()")

	world, err := facade.loadWorld(request.DatabasePath, request.World)
	if err != nil {
		return nil, err
	}

	response, err := world.UpgradeSmartContract()
	if err != nil {
		return nil, err
	}

	dumpResponse(&response)
	return response, err
}

// RunSmartContract -
func (facade *DebugFacade) RunSmartContract(request RunRequest) {
	log.Debug("DebugFacade.RunSmartContract()")
}

// QuerySmartContract -
func (facade *DebugFacade) QuerySmartContract(request QueryRequest) {
	log.Debug("DebugFacade.QuerySmartContracts()")
}

// CreateAccount -
func (facade *DebugFacade) CreateAccount(request CreateAccountRequest) (*CreateAccountResponse, error) {
	log.Debug("DebugFacade.CreateAccount()")

	world, err := facade.loadWorld(request.DatabasePath, request.World)
	if err != nil {
		return nil, err
	}

	response, err := world.CreateAccount(request)
	if err != nil {
		return nil, err
	}

	dumpResponse(&response)
	return response, err
}

func dumpResponse(response interface{}) {
	data, err := json.MarshalIndent(response, "", "\t")
	if err != nil {
		fmt.Println("{}")
	}

	fmt.Println(string(data))
}
