package model

type Device struct {
	ID         int        `json:"id"`
	UUID       string     `json:"uuid"`
	DeviceType DeviceType `json:"device_type"`
	Name       string     `json:"name"`
	State      string     `json:"state"`
}
