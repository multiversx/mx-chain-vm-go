package gasschedules

//go:generate go run scripts/includetoml.go

// GetV3 yields the schedule V3
func GetV3() string {
	return gasScheduleV3
}

// GetV4 yields the schedule V4
func GetV4() string {
	return gasScheduleV4
}
