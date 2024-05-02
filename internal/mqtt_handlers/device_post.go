package mqtt_handlers

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"log"
)

func HandlePost(client MQTT.Client, msg MQTT.Message) {
	log.Printf("Data received from device: %s", msg.Payload())

	// Process the data here (e.g., store in database, perform actions)
}
