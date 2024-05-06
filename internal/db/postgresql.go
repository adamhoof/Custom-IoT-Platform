package db

import (
	"NSI-semester-work/internal/model"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
)

// Database holds the connection pool to the database
type Database struct {
	*sql.DB
}

// NewDatabase creates a new Database connection
func NewDatabase(dataSourceName string) (*Database, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &Database{db}, nil
}

// Disconnect wraps the sql.DB Close method
func (db *Database) Disconnect() error {
	return db.DB.Close()
}

// RegisterDevice registers or authenticates a new device in the database
func (db *Database) RegisterDevice(device *model.Device) error {
	query := `
        INSERT INTO devices (uuid, action_template_id, device_name, custom_actions)
        VALUES ($1, $2, $3, $4) 
        ON CONFLICT (uuid) DO UPDATE SET last_login = NOW()`

	//if Valid -> use String, else use Null
	var customActions sql.NullString
	if device.CustomActions != "" {
		customActions = sql.NullString{String: device.CustomActions, Valid: true}
	}
	_, err := db.Exec(query, device.UUID, device.ActionsTemplateId, device.Name, customActions)
	if err != nil {
		return fmt.Errorf("failed insert device %s\n", err)
	}
	return nil
}

func (db *Database) FetchDeviceNamesAndIds() (devices []model.Device, err error) {
	rows, err := db.Query(`
			SELECT devices.device_id, device_name
			FROM devices`)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {

		}
	}(rows)

	for rows.Next() {
		var device model.Device
		if err = rows.Scan(&device.ID, &device.Name); err != nil {
			return nil, err
		}
		devices = append(devices, device)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return devices, nil
}

func (db *Database) FetchDeviceWithActions(deviceId int) (*model.Device, error) {
	query := `
		SELECT devices.device_id, devices.device_name, action_templates.actions, COALESCE(devices.custom_actions, '{}')
		FROM devices
		JOIN action_templates ON devices.action_template_id = action_templates.action_template_id
		WHERE devices.device_id = $1;
	`

	row := db.QueryRow(query, deviceId)
	var device model.Device
	var templateActions, customActions sql.NullString

	if err := row.Scan(&device.ID, &device.Name, &templateActions, &customActions); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no device found with ID %d", deviceId)
		}
		return nil, err
	}

	device.TemplateActions = templateActions.String
	device.CustomActions = customActions.String

	return &device, nil
}

func (db *Database) CreateDashboard(name string) (dashboardId int, err error) {
	err = db.QueryRow(`INSERT INTO dashboards (name) VALUES ($1) RETURNING dashboard_id`, name).Scan(&dashboardId)
	if err != nil {
		return 0, err
	}
	return dashboardId, nil
}

func (db *Database) InsertDevicesToDashboard(dashboardId int, devices []model.DeviceInDashboard) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	for _, device := range devices {
		shownActionsJSON, err := json.Marshal(device.ShownActions)
		if err != nil {
			tx.Rollback() // Handle rollback in case of error
			return fmt.Errorf("error marshaling shown actions: %v", err)
		}

		_, err = tx.Exec(`INSERT INTO devices_in_dashboard (device_id, dashboard_id, position_in_dashboard, shown_actions) VALUES ($1, $2, $3, $4)`,
			device.Device.ID, dashboardId, device.Position, string(shownActionsJSON))
		if err != nil {
			tx.Rollback() // Handle rollback in case of error
			return fmt.Errorf("error inserting device into dashboard: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}
	return nil
}

func (db *Database) FetchDashboards() ([]model.Dashboard, error) {
	var dashboards []model.Dashboard

	rows, err := db.Query(`SELECT dashboard_id, name FROM dashboards`)
	if err != nil {
		return nil, err // Return nil slice and the error
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	for rows.Next() {
		var d model.Dashboard
		err := rows.Scan(&d.DashboardId, &d.Name)
		if err != nil {
			return nil, err
		}
		dashboards = append(dashboards, d)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return dashboards, nil
}

func (db *Database) FetchDashboardContents(dashboardID int) ([]model.DeviceInDashboard, string, error) {
	var devices []model.DeviceInDashboard
	var dashboardName string

	query := `
    SELECT d.device_id, d.device_name, did.shown_actions, did.position_in_dashboard, dash.name
    FROM devices_in_dashboard did
    JOIN devices d ON did.device_id = d.device_id
    JOIN dashboards dash ON did.dashboard_id = dash.dashboard_id
    WHERE did.dashboard_id = $1
    ORDER BY did.position_in_dashboard
    `
	rows, err := db.Query(query, dashboardID)
	if err != nil {
		return nil, "", fmt.Errorf("error querying database: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var device model.DeviceInDashboard
		var shownActionsJSON string // This holds the JSON string from the database

		if err := rows.Scan(&device.Device.ID, &device.Device.Name, &shownActionsJSON, &device.Position, &dashboardName); err != nil {
			return nil, "", fmt.Errorf("error scanning row: %v", err)
		}

		// Initialize the map to store action name to action type mappings
		device.ShownActions = make(map[string]string)
		// Unmarshal the JSON string into the map
		if err := json.Unmarshal([]byte(shownActionsJSON), &device.ShownActions); err != nil {
			return nil, "", fmt.Errorf("error unmarshaling JSON: %v", err)
		}

		devices = append(devices, device)
	}

	if err = rows.Err(); err != nil {
		return nil, "", fmt.Errorf("error processing rows: %v", err)
	}

	return devices, dashboardName, nil
}

func (db *Database) FetchTemplateActions(deviceType model.DeviceType) (actionTemplateId int, err error) {
	query := `SELECT action_template_id FROM action_templates WHERE device_type = $1;`

	row := db.QueryRow(query, deviceType)
	err = row.Scan(&actionTemplateId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return -1, sql.ErrNoRows
		}
		return -1, err
	}

	return actionTemplateId, nil
}
