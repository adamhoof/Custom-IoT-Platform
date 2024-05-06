package db

import (
	"NSI-semester-work/internal/model"
	"database/sql"
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

func (db *Database) FetchDevicesWithActions() (devices []model.Device, err error) {
	rows, err := db.Query(`
			SELECT devices.device_id, uuid, device_name, actions, coalesce(custom_actions, '{}')
			FROM devices
			JOIN action_templates ON devices.action_template_id = action_templates.action_template_id
		`)
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
		var templateActions, customActions sql.NullString
		if err = rows.Scan(&device.ID, &device.UUID, &device.Name, &templateActions, &customActions); err != nil {
			return nil, err
		}
		device.TemplateActions = templateActions.String
		device.CustomActions = customActions.String
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

// GetDeviceIDByUUID returns the database ID for a given device Id.
func (db *Database) GetDeviceIDByUUID(uuid string) (int, error) {
	var deviceID int
	err := db.QueryRow("SELECT device_id FROM devices WHERE uuid = $1", uuid).Scan(&deviceID)
	if err != nil {
		return 0, err
	}
	return deviceID, nil
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
		_, err := tx.Exec(`INSERT INTO devices_in_dashboard (device_id, dashboard_id, position_in_dashboard,shown_actions) VALUES ($1, $2, $3, $4)`,
			device.Device.ID, dashboardId, device.Position, device.ShownActions)
		if err != nil {
			err = tx.Rollback()
			if err != nil {
				return err
			}
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
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

func (db *Database) FetchDevicesInDashboard(dashboardId int) (devices []model.DeviceInDashboard, err error) {
	rows, err := db.Query(`
			SELECT  devices.device_id, position_in_dashboard, shown_actions, devices.device_name, action_templates.device_type
			FROM devices_in_dashboard join devices on devices.device_id = devices_in_dashboard.device_id join action_templates on devices.action_template_id = action_templates.action_template_id
			WHERE dashboard_id = $1
			ORDER BY position_in_dashboard`, dashboardId)

	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	// Iterate over the rows in the result set
	for rows.Next() {
		var device model.DeviceInDashboard
		err := rows.Scan(&device.Device.ID, &device.Position, &device.ShownActions, &device.Device.Name, &device.Device.DeviceType)
		if err != nil {
			return nil, err
		}
		devices = append(devices, device)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return devices, nil
}

func (db *Database) FetchDashboardName(dashboardId int) (name string, err error) {
	err = db.QueryRow(`SELECT name FROM dashboards where dashboard_id = $1`, dashboardId).Scan(&name)
	if err != nil {
		return "", err
	}
	return name, nil
}

func (db *Database) FetchDashboardContents(dashboardId int) (name string, devices []model.DeviceInDashboard, err error) {
	name, err = db.FetchDashboardName(dashboardId)
	if err != nil {
		fmt.Printf("unable to fetch name %s\n", err)
	}

	devices, err = db.FetchDevicesInDashboard(dashboardId)
	if err != nil {
		fmt.Printf("unable to fetch devices %s\n", err)
		return "", nil, err
	}
	return name, devices, nil
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

/*func (db *Database) GetDeviceStateByID(id int) (state string, err error) {
	err = db.QueryRow("SELECT state FROM devices WHERE device_id = $1", id).Scan(&state)
	if err != nil {
		return "", err
	}
	return state, nil
}*/

/*func (db *Database) UpdateDeviceState(uuid string, state string) error {
	stmt, err := db.Prepare("UPDATE devices SET state = $1 WHERE uuid = $2")
	if err != nil {
		return err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {

		}
	}(stmt)

	// Execute the statement with the JSON payload
	_, err = stmt.Exec(state, uuid)
	if err != nil {
		return err
	}
	return nil
}
*/
