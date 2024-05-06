package model

type DeviceInDashboard struct {
	Device       Device
	ShownActions map[string]string
	Position     int
}
