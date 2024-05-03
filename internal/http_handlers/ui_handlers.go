package http_handlers

import (
	"NSI-semester-work/internal/db"
	"NSI-semester-work/internal/model"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

func HandleHome(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("ui/html/home.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, nil) // Pass nil or any actual data structure if needed
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
	}
}

func DashboardCreatorHandler(w http.ResponseWriter, r *http.Request, database *db.Database) {
	t, err := template.ParseFiles("ui/html/dashboard_creator.html")
	if err != nil {
		fmt.Printf("failed to load dashboard creator template %s\n", err)
		http.Error(w, "Failed to load the dashboard creator template", http.StatusInternalServerError)
		return
	}
	devices, err := database.FetchDevices()
	if err != nil {
		fmt.Println("failed to fetch devices")
		http.Error(w, "Failed to fetch devices", http.StatusInternalServerError)
		return
	}
	// Render template with devices
	if err := t.Execute(w, map[string]interface{}{"Devices": devices}); err != nil {
		fmt.Println(devices)
		fmt.Printf("error executing template %s\n", err)
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}
}

func DeviceFeaturesHandler(w http.ResponseWriter, r *http.Request, database *db.Database) {
	deviceType := ""
	uuid := ""
	query := r.URL.Query()
	if deviceTypeParam, ok := query["deviceType"]; ok {
		deviceType = strings.Split(deviceTypeParam[0], "=")[0]
	}
	if uuidParam, ok := query["uuid"]; ok {
		uuid = strings.Split(uuidParam[0], "=")[0]
	}

	var features []string
	switch deviceType {
	case model.DeviceTypeOnOffDevice.String():
		features = []string{"On", "Off"}
	case model.DeviceTypeSingleMetricSensor.String():
		features = []string{"LastValue"}
	}

	t, err := template.ParseFiles("ui/html/device_features_template.html")
	if err != nil {
		fmt.Printf("error loading feature template %s\n", err)
		http.Error(w, "Error loading feature template", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Features": features,
		"UUID":     uuid,
	}
	if err := t.Execute(w, data); err != nil {
		fmt.Printf("error executing feature template %s\n", err)
		http.Error(w, "Error executing feature template", http.StatusInternalServerError)
		return
	}

}
