package http_handlers

import (
	"NSI-semester-work/internal/db"
	"NSI-semester-work/internal/model"
	"encoding/json"
	"fmt"
	"net/http"
)

func SseStateHandler(w http.ResponseWriter, r *http.Request, db *db.Database, sseChannel chan model.Update) {
	fmt.Println("setting up a new connection")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	ctx := r.Context()

	for {
		select {
		case update := <-sseChannel:
			// Prepare data to send via SSE.
			res := fmt.Sprintf("event: stateUpdate-%d-%s\ndata: %s\n\n", update.DeviceID, update.ActionName, update.State)
			fmt.Println(res) // Logging for debugging purposes.

			// Map to update device state in the database.
			stateMap := map[string]interface{}{
				update.ActionName: update.State,
			}

			// Update state in the database.
			err := db.UpdateDeviceState(update.DeviceID, stateMap)
			if err != nil {
				fmt.Printf("Unable to update state: %s\n", err)
				continue // Decide if you want to skip sending the event.
			}

			// Serialize data into JSON format for sending.
			jsonData, err := json.Marshal(map[string]interface{}{
				"deviceID":   update.DeviceID,
				"actionName": update.ActionName,
				"state":      update.State,
			})
			if err != nil {
				fmt.Println("Failed to serialize update to JSON:", err)
				continue // Skip this update if serialization fails.
			}

			// Send the SSE data to the client.
			_, err = fmt.Fprintf(w, "data: %s\n\n", jsonData)
			if err != nil {
				fmt.Println("Error printing state update data:", err)
				continue
			}
			flusher.Flush()

		case <-ctx.Done():
			// Handle client disconnection.
			fmt.Println("Client disconnected, closing SSE connection.")
			return // Exit the handler when the client disconnects.
		}
	}
}
