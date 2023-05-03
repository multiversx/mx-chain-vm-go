package worldmock

import (
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
)

var _ vmcommon.GuardedAccountHandler = (*GuardedAccountHandlerStub)(nil)

// GuardedAccountHandlerStub -
type GuardedAccountHandlerStub struct {
	GetActiveGuardianCalled    func(handler vmcommon.UserAccountHandler) ([]byte, error)
	SetGuardianCalled          func(uah vmcommon.UserAccountHandler, guardianAddress []byte, txGuardianAddress []byte, guardianServiceUID []byte) error
	CleanOtherThanActiveCalled func(uah vmcommon.UserAccountHandler)
}

// GetActiveGuardian -
func (gahs *GuardedAccountHandlerStub) GetActiveGuardian(handler vmcommon.UserAccountHandler) ([]byte, error) {
	if gahs.GetActiveGuardianCalled != nil {
		return gahs.GetActiveGuardianCalled(handler)
	}
	return nil, nil
}

// SetGuardian -
func (gahs *GuardedAccountHandlerStub) SetGuardian(uah vmcommon.UserAccountHandler, guardianAddress []byte, txGuardianAddress []byte, guardianServiceUID []byte) error {
	if gahs.SetGuardianCalled != nil {
		return gahs.SetGuardianCalled(uah, guardianAddress, txGuardianAddress, guardianServiceUID)
	}
	return nil
}

// CleanOtherThanActive -
func (gahs *GuardedAccountHandlerStub) CleanOtherThanActive(uah vmcommon.UserAccountHandler) {
	if gahs.CleanOtherThanActiveCalled != nil {
		gahs.CleanOtherThanActiveCalled(uah)
	}
}

// IsInterfaceNil -
func (gahs *GuardedAccountHandlerStub) IsInterfaceNil() bool {
	return gahs == nil
}
