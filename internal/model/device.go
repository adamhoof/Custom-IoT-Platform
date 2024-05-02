package model

type Device struct {
	UUID       string     `json:"uuid"`
	DeviceType DeviceType `json:"device_type"`
	Name       string     `json:"name"`
}
