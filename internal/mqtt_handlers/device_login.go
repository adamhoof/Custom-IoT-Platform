package mqtt_handlers

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"log"
	"os"
)

func HandleLogin(client MQTT.Client, msg MQTT.Message) {
	// Assume msg.Payload() contains the device ID and credentials
	log.Printf("Login request from device: %s", msg.Payload())

	// Here you would check the device credentials against your database
	// For now, we assume the login is successful and publish an acknowledgement
	client.Publish(os.Getenv("MQTT_LOGIN_RESPONSE_TOPIC"), 0, false, "Login successful")
}
