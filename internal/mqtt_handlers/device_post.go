package mqtt_handlers

import (
	"NSI-semester-work/internal/db"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"log"
	"strings"
)

func ValueProvidedHandler(msg MQTT.Message, database *db.Database) {
	topic := msg.Topic()
	log.Printf("Message received on topic: %s", topic)

	// Assume topic structure is "provide_value/<device_uuid>"
	parts := strings.Split(topic, "/")
	if len(parts) != 2 {
		log.Println("Invalid topic format")
		return
	}
	uuid := parts[1]

	payloadStr := string(msg.Payload())
	deviceId, err := database.GetDeviceIDByUUID(uuid)
	if err != nil {
		log.Printf("Error retrieving device ID for UUID %s: %s", uuid, err)
		return
	}

	// Insert or update the provided value for the device in the database
	if err := database.InsertProvidedValue(deviceId, payloadStr); err != nil {
		log.Printf("Error updating provided value in the database for device %s: %s", uuid, err)
		return
	}

	// Log successful update
	log.Printf("Updated provided value for device %s with payload: %s", uuid, payloadStr)
}
