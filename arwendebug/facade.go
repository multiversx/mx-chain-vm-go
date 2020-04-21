package arwendebug

import (
	logger "github.com/ElrondNetwork/elrond-go-logger"
)

var log = logger.GetOrCreate("arwendebug")

// DebugFacade -
type DebugFacade struct {
}

// NewDebugFacade -
func NewDebugFacade() {
}

// StartServer -
func (facade *DebugFacade) StartServer(address string) {
	log.Debug("DebugFacade.StartServer()")
	startServer(address)
}

// DeploySmartContract -
func (facade *DebugFacade) DeploySmartContract(request DeployRequest) (DeployResponse, error) {
	log.Debug("DebugFacade.DeploySmartContract()")

	session, err := facade.loadSession(request.DatabasePath, request.Session)
	if err != nil {
		return DeployResponse{}, err
	}

	return session.DeploySmartContract(request)
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

// UpgradeSmartContract -
func (facade *DebugFacade) UpgradeSmartContract(request UpgradeRequest) error {
	log.Debug("DebugFacade.UpgradeSmartContract()")

	session, err := facade.loadSession(request.DatabasePath, request.Session)
	if err != nil {
		return err
	}

	return session.UpgradeSmartContract()
}

// RunSmartContract -
func (facade *DebugFacade) RunSmartContract(request RunRequest) {
	log.Debug("DebugFacade.RunSmartContract()")
}

// QuerySmartContract -
func (facade *DebugFacade) QuerySmartContract(query QueryRequest) {
	log.Debug("DebugFacade.QuerySmartContracts()")
}
