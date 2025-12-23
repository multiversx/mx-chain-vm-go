package config

type gasScheduleFactory struct{}

func NewGasScheduleFactory() *gasScheduleFactory {
	return &gasScheduleFactory{}
}

func (f *gasScheduleFactory) CreateGasSchedule(gasMap GasScheduleMap) (GasSchedule, error) {
	return CreateGasConfig(gasMap)
}

func (f *gasScheduleFactory) IsInterfaceNil() bool {
	return f == nil
}
