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
