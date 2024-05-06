package model

type DeviceInDashboard struct {
	Device       Device
	ShownActions map[string]string
	Position     int
}

/*func (d *DeviceInDashboard) FunctionalitiesList() ([]string, error) {
	var functionalities []string
	err := json.Unmarshal([]byte(d.ShownActions), &functionalities)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling functionalities: %w", err)
	}
	return functionalities, nil
}

// HasFunctionality checks if a device has a specific functionality.
func (d *DeviceInDashboard) HasFunctionality(funcName string) bool {
	functionalities, err := d.FunctionalitiesList()
	if err != nil {
		fmt.Println("Error parsing functionalities:", err)
		return false
	}
	for _, funct := range functionalities {
		if funct == funcName {
			return true
		}
	}
	return false
}*/
