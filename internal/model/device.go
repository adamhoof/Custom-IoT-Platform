package model

type Device struct {
	ID                int        `json:"id"`
	UUID              string     `json:"uuid"`
	Name              string     `json:"name"`
	TemplateActions   string     `json:"template_actions"`
	CustomActions     string     `json:"custom_actions"`
	DeviceType        DeviceType `json:"device_type"`
	ActionsTemplateId int        `json:"actions_template_id"`
}
