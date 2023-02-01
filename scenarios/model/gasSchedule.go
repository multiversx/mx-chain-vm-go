package scenjsonmodel

// GasSchedule encodes the gas model to be used in scenario tests
type GasSchedule int

const (
	// GasScheduleDefault indicates that the scenario should use whatever the default gas model is.
	// Should be the latest version of the mainnet gas schedule.
	GasScheduleDefault GasSchedule = iota

	// GasScheduleDummy is a dummy model, with all costs set to 1.
	GasScheduleDummy

	// GasScheduleV3 is currently used on mainnet.
	GasScheduleV3

	// GasScheduleV4 is currently used on mainnet.
	GasScheduleV4
)
