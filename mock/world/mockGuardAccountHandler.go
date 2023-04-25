package worldmock

import vmcommon "github.com/multiversx/mx-chain-vm-common-go"

// MockGuardedAccountHandler -
type MockGuardedAccountHandler struct{}

// NewMockGuardedAccountHandler -
func NewMockGuardedAccountHandler() *MockGuardedAccountHandler {
	return &MockGuardedAccountHandler{}
}

// GetActiveGuardian -
func (mah *MockGuardedAccountHandler) GetActiveGuardian(_ vmcommon.UserAccountHandler) ([]byte, error) {
	return nil, nil
}

// SetGuardian -
func (mah *MockGuardedAccountHandler) SetGuardian(_ vmcommon.UserAccountHandler, _ []byte, _ []byte, _ []byte) error {
	return nil
}

// CleanOtherThanActive -
func (mah *MockGuardedAccountHandler) CleanOtherThanActive(_ vmcommon.UserAccountHandler) {
}

// IsInterfaceNil -
func (mah *MockGuardedAccountHandler) IsInterfaceNil() bool {
	return mah == nil
}
