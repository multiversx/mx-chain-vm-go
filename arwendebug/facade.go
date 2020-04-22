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

	session, err := facade.loadSession(request.DatabasePath, request.Session)
	if err != nil {
		return nil, err
	}

	response, err := session.DeploySmartContract(request)
	if err != nil {
		return nil, err
	}

	dumpResponse(&response)
	return response, err
}

func (facade *DebugFacade) loadSession(databaseRootPath string, sessionID string) (*session, error) {
	database := facade.loadDatabase(databaseRootPath)
	return database.loadSession(sessionID)
}

func (facade *DebugFacade) loadDatabase(rootPath string) *database {
	// TODO: use factory
	database := NewDatabase(rootPath)
	return database
}

func (facade *DebugFacade) saveSession() {

}

// UpgradeSmartContract -
func (facade *DebugFacade) UpgradeSmartContract(request UpgradeRequest) (*UpgradeResponse, error) {
	log.Debug("DebugFacade.UpgradeSmartContract()")

	session, err := facade.loadSession(request.DatabasePath, request.Session)
	if err != nil {
		return nil, err
	}

	response, err := session.UpgradeSmartContract()
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

	session, err := facade.loadSession(request.DatabasePath, request.Session)
	if err != nil {
		return nil, err
	}

	response, err := session.CreateAccount(request)
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
