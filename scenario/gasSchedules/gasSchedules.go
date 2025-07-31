package gasschedules

import _ "embed"

//go:embed gasScheduleV3.toml
var gasScheduleV3 string

// GetV3 yields the schedule V3
func GetV3() string {
	return gasScheduleV3
}

//go:embed gasScheduleV4.toml
var gasScheduleV4 string

// GetV4 yields the schedule V4
func GetV4() string {
	return gasScheduleV4
}
