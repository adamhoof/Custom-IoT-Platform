``` sql
-- Create a custom type for device types if you have a known set of types
CREATE TYPE device_type AS ENUM ('temperature_sensor', 'humidity_sensor', 'soil_moisture_sensor', 'on_off_device');

-- Table for storing action templates with a JSONB validation check
CREATE TABLE action_templates
(
    action_template_id SERIAL PRIMARY KEY,
    device_type        device_type NOT NULL UNIQUE,
    actions            JSONB NOT NULL CHECK (
        (actions ? 'Temperature' AND jsonb_typeof(actions -> 'Temperature') = 'string') OR
        (actions ? 'Humidity' AND jsonb_typeof(actions -> 'Humidity') = 'string') OR
        (actions ? 'Soil_moisture' AND jsonb_typeof(actions -> 'Soil_moisture') = 'string') OR
        (actions ? 'On' AND jsonb_typeof(actions -> 'On') = 'string') OR
        (actions ? 'Off' AND jsonb_typeof(actions -> 'Off') = 'string')
    )
);

-- Populate the action_templates table
INSERT INTO action_templates (device_type, actions) VALUES
('temperature_sensor', '{"Temperature": "provide_value", "Interval_ms": "numberInput"}'::jsonb),
('humidity_sensor', '{"Humidity": "provide_value", "Interval_ms": "numberInput"}'::jsonb),
('soil_moisture_sensor', '{"Soil_moisture": "provide_value", "Interval_ms": "numberInput"}'::jsonb),
('on_off_device', '{"On": "command", "Off": "command"}'::jsonb);

-- Table for storing device information
CREATE TABLE devices
(
    device_id          SERIAL PRIMARY KEY,
    uuid               UUID UNIQUE NOT NULL,
    device_name        TEXT NOT NULL,
    action_template_id INTEGER REFERENCES action_templates (action_template_id),
    custom_actions     JSONB,
    last_login         TIMESTAMP(0) DEFAULT CURRENT_TIMESTAMP
);

-- Table for storing dashboard information
CREATE TABLE dashboards
(
    dashboard_id SERIAL PRIMARY KEY,
    name         TEXT UNIQUE NOT NULL
);

-- Table for storing device to dashboard mappings with position checking
CREATE TABLE devices_in_dashboard
(
    device_id             INT NOT NULL,
    dashboard_id          INT NOT NULL,
    position_in_dashboard INT DEFAULT -1 CHECK (position_in_dashboard >= -1),
    shown_actions         JSONB,
    PRIMARY KEY (device_id, dashboard_id),
    FOREIGN KEY (device_id) REFERENCES devices (device_id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (dashboard_id) REFERENCES dashboards (dashboard_id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- Table for storing sensor data as a hypertable
CREATE TABLE sensor_data
(
    timestamp TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    device_id INT NOT NULL,
    data      JSONB NOT NULL,
    CONSTRAINT fk_device FOREIGN KEY (device_id) REFERENCES devices (device_id)
);

-- Convert sensor_data table to a hypertable
SELECT create_hypertable('sensor_data', 'timestamp');

-- Indexes for optimized query performance
CREATE INDEX idx_timestamp_device_id ON sensor_data (timestamp DESC, device_id);
CREATE INDEX idx_device_id ON sensor_data (device_id);
CREATE INDEX idx_timestamp ON sensor_data (timestamp DESC);
CREATE INDEX idx_device_type ON action_templates (device_type);

```