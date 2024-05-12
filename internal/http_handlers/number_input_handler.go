package http_handlers

import (
	"NSI-semester-work/internal/db"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"net/http"
	"strconv"
)

func NumberInputHandler(w http.ResponseWriter, r *http.Request, database *db.Database, mqttClient MQTT.Client) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	deviceIdStr := r.FormValue("deviceID")
	actionName := r.FormValue("actionName")
	inputValue := r.FormValue("inputValue")

	deviceId, err := strconv.Atoi(deviceIdStr)
	if err != nil {
		http.Error(w, "Invalid device ID", http.StatusBadRequest)
		return
	}

	deviceUuid, err := database.GetDeviceUUID(deviceId)

	topic := fmt.Sprintf("number_input/%s/%s", deviceUuid, actionName)
	if token := mqttClient.Publish(topic, 0, false, inputValue); token.Wait() && token.Error() != nil {
		http.Error(w, "Failed to send command", http.StatusInternalServerError)
		return
	}
}
