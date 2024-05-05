package db

import (
	"NSI-semester-work/internal/model"
	"database/sql"
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
	var deviceTypeID int
	err := db.QueryRow("SELECT device_type_id FROM device_types WHERE name = $1", device.DeviceType).Scan(&deviceTypeID)
	if err != nil {
		return fmt.Errorf("failed to search for device type id %s\n", err)
	}

	query := `
        INSERT INTO devices (uuid, device_type_id, name)
        VALUES ($1, $2, $3) 
        ON CONFLICT (uuid) DO UPDATE SET last_login = NOW()
    `
	_, err = db.Exec(query, device.UUID, deviceTypeID, device.Name)
	if err != nil {
		return fmt.Errorf("failed insert device %s\n", err)
	}
	return nil
}

func (db *Database) FetchDevices() (devices []model.Device, err error) {
	rows, err := db.Query(`
        SELECT d.uuid, dt.name AS device_type, d.name
        FROM devices d
        JOIN device_types dt ON d.device_type_id = dt.device_type_id
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
		if err = rows.Scan(&device.UUID, &device.DeviceType, &device.Name); err != nil {
			return nil, err
		}

		devices = append(devices, device)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return devices, nil
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

func (db *Database) CreateDashboard(name string) (int, error) {
	var dashboardID int
	// Query to insert the dashboard and return the new ID
	err := db.QueryRow(`INSERT INTO dashboard (name) VALUES ($1) RETURNING dashboard_id`, name).Scan(&dashboardID)
	if err != nil {
		return 0, err
	}
	return dashboardID, nil
}

func (db *Database) InsertDevicesToDashboard(dashboardID int, devices []model.DeviceInDashboard) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	for _, device := range devices {
		_, err := tx.Exec(`INSERT INTO device_in_dashboard (device_id, dashboard_id, position_in_dashboard, shown_functionalities) VALUES ($1, $2, $3, $4)`,
			device.Device.ID, dashboardID, device.Position, device.Functionalities)
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

	rows, err := db.Query(`SELECT dashboard_id, name FROM dashboard`)
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
			SELECT devices.device_id, position_in_dashboard, shown_functionalities, devices.name, device_types.name, devices.state
			FROM device_in_dashboard join devices on devices.device_id = device_in_dashboard.device_id join device_types on devices.device_type_id = device_types.device_type_id
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
		err := rows.Scan(&device.Device.ID, &device.Position, &device.Functionalities, &device.Device.Name, &device.Device.DeviceType, &device.Device.State)
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
	err = db.QueryRow(`SELECT name FROM dashboard where dashboard_id = $1`, dashboardId).Scan(&name)
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

func (db *Database) GetDeviceStateByID(id int) (state string, err error) {
	err = db.QueryRow("SELECT state FROM devices WHERE device_id = $1", id).Scan(&state)
	if err != nil {
		return "", err
	}
	return state, nil
}

func (db *Database) UpdateDeviceState(uuid string, state string) error {
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
