package mqtt_handlers

import (
	"NSI-semester-work/internal/db"
	"NSI-semester-work/internal/model"
	"encoding/json"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"log"
	"os"
)

func HandleDeviceLogin(client MQTT.Client, msg MQTT.Message, database *db.Database) {
	var device model.Device
	if err := json.Unmarshal(msg.Payload(), &device); err != nil {
		log.Printf("Error decoding JSON: %s", err)
		return
	}
	log.Println(device)

	err := database.RegisterDevice(&device)
	if err != nil {
		log.Println(err)
	}

	token := client.Publish(os.Getenv("MQTT_LOGIN_RESPONSE_TOPIC"), 0, false, "Success")
	token.Wait()
}
