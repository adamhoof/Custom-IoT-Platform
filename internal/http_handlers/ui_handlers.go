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

func CreateDashboard(w http.ResponseWriter, r *http.Request, db *db.Database) {
	// Parse the request form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}
	fmt.Println("Form data:", r.Form)

	// Extract and validate dashboard name
	dashboardName := r.FormValue("dashboardName")
	if dashboardName == "" {
		http.Error(w, "Dashboard name is required", http.StatusBadRequest)
		return
	}

	// Remove the dashboardName key from the form to process device entries
	r.Form.Del("dashboardName")

	var deviceEntries []model.DeviceInDashboard
	position := 0
	for key, values := range r.Form {
		// Encode functionalities into a JSON string
		functionalitiesJSON, err := json.Marshal(values)
		if err != nil {
			http.Error(w, "Error encoding functionalities", http.StatusInternalServerError)
			return
		}

		deviceID, err := db.GetDeviceIDByUUID(key)
		deviceEntries = append(deviceEntries, model.DeviceInDashboard{
			Id:              deviceID,
			Functionalities: string(functionalitiesJSON),
			Position:        position,
		})
		position++
	}

	dashboardID, err := db.CreateDashboard(dashboardName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create dashboard: %v", err), http.StatusInternalServerError)
		return
	}
	if err = db.InsertDevicesToDashboard(dashboardID, deviceEntries); err != nil {
		http.Error(w, fmt.Sprintf("Failed to save devices: %v", err), http.StatusInternalServerError)
		return
	}

	// Redirect or update frontend as needed
	fmt.Fprintln(w, "Dashboard saved successfully!") // Adjust based on your frontend needs
}
