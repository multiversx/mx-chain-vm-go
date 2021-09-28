package mandosjsonmodel

// GasSchedule encodes the gas model to be used in mandos tests
type GasSchedule int

const (
	// GasScheduleDefault indicates that the mandos scenario should use whatever the default gas model is.
	// Should be the latest version of the mainnet gas schedule.
	GasScheduleDefault GasSchedule = iota

	// GasScheduleDummy is a dummy model, with all costs set to 1.
	GasScheduleDummy

	// GasScheduleV1 was previously used on mainnet.
	GasScheduleV1

	// GasScheduleV2 was previously used on mainnet.
	GasScheduleV2

	// GasScheduleV3 is currently used on mainnet.
	GasScheduleV3
)
