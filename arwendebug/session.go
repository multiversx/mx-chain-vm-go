package arwendebug

import "github.com/ElrondNetwork/arwen-wasm-vm/mock"

type session struct {
	blockchainHook mock.BlockchainHookMock
}

// NewSession -
func NewSession(blockchainHook *mock.BlockchainHookMock) *session {
	return &session{
		blockchainHook: *blockchainHook,
	}
}

// DeploySmartContract -
func (session *session) DeploySmartContract() error {
	return nil
}

// UpgradeSmartContract -
func (session *session) UpgradeSmartContract() error {
	return nil
}

// RunSmartContract -
func (session *session) RunSmartContract() error {
	return nil
}

// QuerySmartContract -
func (session *session) QuerySmartContract() error {
	return nil
}
