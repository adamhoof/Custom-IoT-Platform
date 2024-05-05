package model

type Dashboard struct {
	DashboardId int
	Name        string
	Devices     []DeviceInDashboard
}
