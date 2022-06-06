package worldmock

// EnableRoundsHandlerMock -
type EnableRoundsHandlerMock struct {
	IsCheckValueOnExecByCallerEnabledValue bool
}

// IsCheckValueOnExecByCallerEnabled -
func (mock *EnableRoundsHandlerMock) IsCheckValueOnExecByCallerEnabled() bool {
	return mock.IsCheckValueOnExecByCallerEnabledValue
}

// IsInterfaceNil -
func (mock *EnableRoundsHandlerMock) IsInterfaceNil() bool {
	return mock == nil
}
