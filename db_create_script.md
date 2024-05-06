``` sql
CREATE TABLE action_templates
(
action_template_id SERIAL PRIMARY KEY,
device_type        TEXT  NOT NULL UNIQUE,
actions            JSONB NOT NULL
);

INSERT INTO action_templates (device_type, actions)
VALUES ('temperature_sensor', '{
"Temperature": "provide_value",
"Interval_ms": "numberInput"
}');

INSERT INTO action_templates (device_type, actions)
VALUES ('humidity_sensor', '{
"Humidity": "provide_value",
"Interval_ms": "numberInput"
}');

INSERT INTO action_templates (device_type, actions)
VALUES ('soil_moisture_sensor', '{
"Soil_moisture": "provide_value",
"Interval_ms": "numberInput"
}');

INSERT INTO action_templates (device_type, actions)
VALUES ('on_off_device', '{
"On": "command",
"Off": "command"
}');

CREATE TABLE devices
(
device_id          SERIAL PRIMARY KEY,
uuid               UUID UNIQUE NOT NULL,
device_name        TEXT        NOT NULL,
action_template_id INTEGER REFERENCES action_templates (action_template_id),
custom_actions     JSONB,
last_login         TIMESTAMP(0) DEFAULT CURRENT_TIMESTAMP
);

create table dashboards
(
dashboard_id SERIAL PRIMARY KEY,
name         text UNIQUE NOT NULL
);


CREATE TABLE devices_in_dashboard
(
device_id             INT NOT NULL,
dashboard_id          INT NOT NULL,
position_in_dashboard INT DEFAULT -1,
shown_actions         jsonb,
PRIMARY KEY (device_id, dashboard_id),
FOREIGN KEY (device_id) REFERENCES devices (device_id)
ON DELETE CASCADE
ON UPDATE CASCADE,
FOREIGN KEY (dashboard_id) REFERENCES dashboards (dashboard_id)
ON DELETE CASCADE
ON UPDATE CASCADE
);


create table sensor_data
(
timestamp TIMESTAMP(0) DEFAULT CURRENT_TIMESTAMP not null,
device_id int references devices (device_id)     not null,
data      jsonb                                  not null

);
```