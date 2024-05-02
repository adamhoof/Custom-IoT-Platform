#include <WiFi.h>
#include <PubSubClient.h>
#include <ArduinoJson.h>
#include "env.h"

WiFiClient espClient;
PubSubClient mqttClient(espClient);

void setup_wifi() {
    delay(10);
    Serial.print("Connecting to ");
    Serial.println(ssid);

    WiFi.begin(ssid, password);
    while (WiFiClass::status() != WL_CONNECTED) {
        delay(500);
        Serial.print(".");
    }

    Serial.println("");
    Serial.println("WiFi connected");
    Serial.println("IP address: ");
    Serial.println(WiFi.localIP());
}

void reconnect_mqtt_mqttClient() {
    while (!mqttClient.connected()) {
        Serial.print("Attempting MQTT connection...");
        if (mqttClient.connect(mqttClientId)) {
            Serial.println("connected");

            StaticJsonDocument<200> doc;
            doc["uuid"] = mqttClientId;
            doc["name"] = name;
            doc["device_type"] = deviceType;

            char jsonBuffer[512];
            serializeJson(doc, jsonBuffer);

            mqttClient.publish(login_request_topic, jsonBuffer);
        } else {
            Serial.print("failed, rc=");
            Serial.print(mqttClient.state());
            Serial.println(" try again in 5 seconds");
            delay(5000);
        }
    }
}

void setup() {
    Serial.begin(115200);
    setup_wifi();
    mqttClient.setServer(mqtt_broker, mqtt_port);

    reconnect_mqtt_mqttClient();
}

void loop() {
    if (!mqttClient.connected()) {
        reconnect_mqtt_mqttClient();
    }
    mqttClient.loop();
}
