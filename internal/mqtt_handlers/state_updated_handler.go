package mqtt_handlers

import (
	"NSI-semester-work/internal/db"
	"NSI-semester-work/internal/model"
	"encoding/json"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"log"
	"strings"
)

func parseMessage(msg MQTT.Message) (uuid string, actionName string, stateValue string) {
	topicParts := strings.Split(msg.Topic(), "/")
	if len(topicParts) < 2 {
		log.Println("Invalid topic format")
		return "", "", ""
	}
	uuid = topicParts[1]

	var jsonData map[string]interface{}
	err := json.Unmarshal(msg.Payload(), &jsonData)
	if err != nil {
		log.Println("JSON parsing error:", err)
		return "", "", ""
	}

	actionName, ok := jsonData["Action_name"].(string)
	if !ok {
		log.Println("Action_name not found or invalid")
		return "", "", ""
	}

	stateValue, ok = jsonData[actionName].(string)
	if !ok {
		log.Printf("State value not found or invalid: %s\n", stateValue)
		return "", "", ""
	}

	return uuid, actionName, stateValue
}
func StateUpdatedHandler(message MQTT.Message, database *db.Database, sseChannel chan model.Update) {
	var update model.Update

	var deviceUuid string
	deviceUuid, update.ActionName, update.State = parseMessage(message)
	deviceId, err := database.GetDeviceIDByUUID(deviceUuid)
	if err != nil {
		fmt.Printf("no such device with this uuid %d\n", update.DeviceID)
		return
	}
	update.DeviceID = deviceId
	sseChannel <- update
}
