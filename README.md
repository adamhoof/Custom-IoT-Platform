IoT Dashboard Application

This application allows users to create and manage custom dashboards featuring various IoT devices. Devices communicate via MQTT and can be dynamically controlled and monitored through a user-friendly web interface.
Features

    Custom Dashboards: Users can create personalized dashboards that aggregate controls and information from multiple devices.
    Device Management: Users can create an arbitrary number of devices that connect to the application, provided they adhere to a specified communication protocol.
    Action Types: Devices can perform various actions such as toggle, number input, provide value, and execute commands based on their capabilities.
    Dynamic Device Interaction: Devices can post state updates and respond to commands dynamically via MQTT topics.
    Template-Driven Configuration: Device types are pre-configured with applicable actions to simplify setup.

Device Protocol

Each device integrated into the system must support the following MQTT topics and functionalities (note that all "uuid" parameters are generated on the end device side):

    login/request/uuid: Devices post to this topic to request login.
    login/response/uuid: Devices subscribe to this topic to receive login confirmation and initial configuration, like receiving their  last available state, before they went offline for example.
    state/uuid: Devices post their current state updates to this topic, for example when a device is commanded to do something, it posts to this topic notifying it executed the command properly.

Action Types

Devices can support these types of actions:

    toggle: Change the state between two conditions.
    number_input: Adjust a setting using numeric input.
    provide_value: Provide data readings like temperature or humidity.
    command: Execute a specific command or action.

Device Types and Templates

The system supports several pre-configured templates that define standard actions for common types of IoT devices:

    Light Switch
        Actions: {"Light_state": "toggle"}
    Temperature Sensor
        Actions: {"Interval_ms": "number_input", "Temperature": "provide_value"}
    Soil Moisture Sensor
        Actions: {"Interval_ms": "number_input", "Soil_moisture": "provide_value"}
    Humidity Sensor
        Actions: {"Humidity": "provide_value", "Interval_ms": "number_input"}
    
Example of connecting to the system:

 Serial.print("Attempting MQTT connection...");
        if (mqttClient.connect(mqttClientId)) {
            Serial.println("connected");
            StaticJsonDocument<200> loginDoc;
            loginDoc["uuid"] = "here goes your generated uuid";
            loginDoc["name"] = "light_switch_bedroom";
            loginDoc["device_type"] = "light_switch";
            loginDoc["custom_actions"] = "json string of custom actions, can look like this or can be empty -> 
            {"action_name_1": "toggle", "action_name_2": "number_input"}"

            char loginJsonBuffer[512];
            serializeJson(loginDoc, loginJsonBuffer);

            mqttClient.subscribe(login_response_topic.c_str());
            mqttClient.publish(login_request_topic.c_str(), loginJsonBuffer);

Setup and Configuration

    Run PostgreSQL+Timescale DB, Mosquitto broker and web server with gui by executing this command in the root directory: 
    docker compose up
