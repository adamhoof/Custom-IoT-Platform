#include <WiFi.h>
#include <PubSubClient.h>
#include <ArduinoJson.h>
#include "env.h"

WiFiClient espClient;
PubSubClient mqttClient(espClient);

unsigned long lastMessageTime = 0;
const long messageInterval = 2000;

void setup_wifi()
{
    Serial.print("Connecting to ");
    Serial.println(ssid);
    WiFi.begin(ssid, password);

    while (WiFiClass::status() != WL_CONNECTED) {
        delay(500);
        Serial.print(".");
    }
    Serial.println("\nWiFi connected");
    Serial.println("IP address: ");
    Serial.println(WiFi.localIP());
}

void mqttCallback(char* topic, byte* payload, unsigned int length)
{
    Serial.print("Message arrived on topic: ");
    Serial.println(topic);
    Serial.print("Message:");

    // Create a character array to store the incoming payload data
    char msg[length + 1];
    strncpy(msg, (char*) payload, length);
    msg[length] = '\0';  // Null-terminate the array
    Serial.println(msg);

    char expectedTopic[100];
    sprintf(expectedTopic, "%s/%s", mqttClientId, "state");

    if (strcmp(topic, expectedTopic) == 0) {
        if (strcmp(msg, "LastValue") == 0) {
            Serial.println("Turning sending last value");
    }else if (strcmp(topic,login_response_topic) == 0){
            char t[50];
            sprintf(t, "%s/%s", mqttClientId, "state");
            mqttClient.subscribe(t);
        }}
}

void reconnect_mqtt()
{
    while (!mqttClient.connected()) {
        Serial.print("Attempting MQTT connection...");
        if (mqttClient.connect(mqttClientId)) {
            Serial.println("connected");
            StaticJsonDocument<200> loginDoc;
            loginDoc["uuid"] = mqttClientId;
            loginDoc["name"] = name;
            loginDoc["device_type"] = deviceType;

            char loginJsonBuffer[512];
            serializeJson(loginDoc, loginJsonBuffer);

            mqttClient.publish(login_request_topic, loginJsonBuffer);
        } else {
            Serial.print("failed, rc=");
            Serial.print(mqttClient.state());
            Serial.println(" try again in 5 seconds");
            delay(5000);
        }
    }
}

void setup()
{
    Serial.begin(115200);
    setup_wifi();
    mqttClient.setServer(mqtt_broker, mqtt_port);
    mqttClient.setCallback(mqttCallback);
    reconnect_mqtt();
}

void sendRandomValue()
{
    if (millis() - lastMessageTime > messageInterval) {
        lastMessageTime = millis();

        double randomValue = random(0, 100) / 1.0;
        char topic[50];
        sprintf(topic, "%s/%s", mqttClientId, "state");

        StaticJsonDocument<200> stateDoc;
        stateDoc["temperature"] = randomValue;

        char stateJsonBuffer[512];
        serializeJson(stateDoc, stateJsonBuffer);

        while (!mqttClient.publish(topic, stateJsonBuffer)) {
            delay(2000);
            mqttClient.publish(topic, stateJsonBuffer);
        }

        Serial.print("Published random value: ");
        Serial.println(randomValue);
    }
}

void loop()
{
    if (!mqttClient.connected()) {
        reconnect_mqtt();
    }
    mqttClient.loop();

    sendRandomValue();
}
