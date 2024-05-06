package model

type DeviceType string

const (
	DeviceTypeOnOffDevice        DeviceType = "on_off"
	DeviceTypeTemperatureSensor  DeviceType = "temperature_sensor"
	DeviceTypeHumiditySensor     DeviceType = "humidity_sensor"
	DeviceTypeSoilMoistureSensor DeviceType = "soil_moisture_sensor"
)

func (dt DeviceType) String() string {
	return string(dt)
}
