package arwendebug

import (
	logger "github.com/ElrondNetwork/elrond-go-logger"
)

var log = logger.GetOrCreate("arwendebug/facade")

// DebugFacade -
type DebugFacade struct {
}

// StartServer -
func (facade *DebugFacade) StartServer(address string) {
	log.Debug("DebugFacade.StartServer()")
	startServer(address)
}

// DeploySmartContract -
func (facade *DebugFacade) DeploySmartContract() {

}

// UpgradeSmartContract -
func (facade *DebugFacade) UpgradeSmartContract() {

}

// RunSmartContract -
func (facade *DebugFacade) RunSmartContract() {

}

// QuerySmartContract -
func (facade *DebugFacade) QuerySmartContract() {

}
