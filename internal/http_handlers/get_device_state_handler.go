package http_handlers

import (
	"NSI-semester-work/internal/db"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

func GetDeviceState(w http.ResponseWriter, r *http.Request, database *db.Database) {
	deviceIdStr := r.PathValue("device_id")
	actionName := r.PathValue("action_name")

	deviceId, err := strconv.Atoi(deviceIdStr)
	if err != nil {
		http.Error(w, "Invalid device ID", http.StatusBadRequest)
		return
	}

	state, err := database.GetDeviceState(deviceId, actionName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Device state not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	_, err = fmt.Fprintln(w, state)
	if err != nil {
		fmt.Printf("unable to print state: %s\n", err)
		return
	}
}
