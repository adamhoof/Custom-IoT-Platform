package http_handlers

import (
	"NSI-semester-work/internal/db"
	"NSI-semester-work/internal/model"
	"encoding/json"
	"fmt"
	"net/http"
)

func SseStateHandler(w http.ResponseWriter, db *db.Database, sseChannel chan model.Update) {
	w.Header().Set("Access-Control-Allow-Origin", "*") // or a more specific domain
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for update := range sseChannel {
		res := fmt.Sprintf("event: stateUpdate-%d-%s\ndata: %s\n\n", update.DeviceID, update.ActionName, update.State)
		fmt.Println(res)
		data := map[string]interface{}{
			"deviceID":   update.DeviceID,
			"actionName": update.ActionName,
			"state":      update.State,
		}

		stateMap := map[string]interface{}{
			update.ActionName: update.State}
		err := db.UpdateDeviceState(update.DeviceID, stateMap)
		if err != nil {
			fmt.Printf("unable to update state: %s\n", err)
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			fmt.Println("Failed to serialize update to JSON:", err)
			continue // Skip this update or handle error appropriately
		}

		_, err = fmt.Fprintf(w, "data: %s\n\n", jsonData)
		if err != nil {
			fmt.Println("error printing state update data!")
		}
		w.(http.Flusher).Flush()
	}
}
