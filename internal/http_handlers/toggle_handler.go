package http_handlers

import (
	"NSI-semester-work/internal/db"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"net/http"
	"strconv"
)

func ToggleHandler(w http.ResponseWriter, r *http.Request, database *db.Database, mqttClient MQTT.Client) {
	deviceIdStr := r.PathValue("device_id")
	actionName := r.PathValue("action_name")

	deviceId, err := strconv.Atoi(deviceIdStr)
	if err != nil {
		http.Error(w, "Invalid device ID", http.StatusBadRequest)
		return
	}

	deviceUuid, err := database.GetDeviceUUID(deviceId)

	topic := fmt.Sprintf("toggle/%s", deviceUuid)
	if token := mqttClient.Publish(topic, 0, false, actionName); token.Wait() && token.Error() != nil {
		http.Error(w, "Failed to send command", http.StatusInternalServerError)
		return
	}
}
