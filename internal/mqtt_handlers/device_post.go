package mqtt_handlers

import (
	"NSI-semester-work/internal/db"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"log"
	"strings"
)

func GetDeviceStateHandler(msg MQTT.Message, database *db.Database) {
	topic := msg.Topic()
	log.Printf("Message received on topic: %s", topic)

	if strings.Contains(topic, "/state") {
		parts := strings.Split(topic, "/state")
		if len(parts) < 1 {
			log.Println("Invalid topic format")
			return
		}
		uuid := parts[0]

		payloadStr := string(msg.Payload())

		if err := database.UpdateDeviceState(uuid, payloadStr); err != nil {
			log.Printf("Error updating device state in DB: %s", err)
			return
		}

		log.Printf("Updated state for device %s with JSON payload", uuid)
	}
}
