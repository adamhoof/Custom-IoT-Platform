package model

type DeviceType string

const (
	DeviceTypeOnOffDevice        DeviceType = "on_off"
	DeviceTypeSingleMetricSensor DeviceType = "single_metric_sensor"
)

func (dt DeviceType) String() string {
	return string(dt)
}
