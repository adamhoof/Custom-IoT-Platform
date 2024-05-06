package http_handlers

import (
	"NSI-semester-work/internal/db"
	"NSI-semester-work/internal/model"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func HandleHome(w http.ResponseWriter, database *db.Database) {
	t, err := template.ParseFiles("ui/html/home.gohtml", "ui/html/dashboard_list.gohtml")
	if err != nil {
		fmt.Printf("error loading template %s\n", err)
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	dashboards, err := database.FetchDashboards()
	if err != nil {
		fmt.Printf("failed to fetch dashboards %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, map[string]interface{}{
		"Dashboards": dashboards,
	})
	if err != nil {
		fmt.Printf("failed to execute template %s\n", err)
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}

}

func DashboardCreatorHandler(w http.ResponseWriter, database *db.Database) {
	t, err := template.ParseFiles("ui/html/dashboard_creator.gohtml")
	if err != nil {
		fmt.Printf("failed to load dashboard creator template %s\n", err)
		http.Error(w, "Failed to load the dashboard creator template", http.StatusInternalServerError)
		return
	}
	devices, err := database.FetchDeviceNamesAndIds()
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

func parseJSONActions(templateActionsStr, customActionsStr string) (map[string]interface{}, map[string]interface{}, error) {
	var templateActions, customActions map[string]interface{}
	if err := json.Unmarshal([]byte(templateActionsStr), &templateActions); err != nil {
		return nil, nil, err
	}
	if err := json.Unmarshal([]byte(customActionsStr), &customActions); err != nil {
		return nil, nil, err
	}
	return templateActions, customActions, nil
}

func DeviceFeaturesHandler(w http.ResponseWriter, r *http.Request, database *db.Database) {
	stringId := r.PathValue("id")
	if stringId == "" {
		http.Error(w, "Device ID parameter is missing", http.StatusBadRequest)
		return
	}

	deviceId, err := strconv.Atoi(stringId)
	if err != nil {
		http.Error(w, "Invalid device ID", http.StatusBadRequest)
		return
	}

	device, err := database.FetchDeviceWithActions(deviceId) // Make sure this function name matches the actual function
	if err != nil {
		log.Printf("error fetching device with actions: %s", err)
		http.Error(w, "Failed to fetch device details", http.StatusInternalServerError)
		return
	}

	t, err := template.ParseFiles("ui/html/device_features_template.gohtml")
	if err != nil {
		log.Printf("error loading feature template: %s", err)
		http.Error(w, "Error loading feature template", http.StatusInternalServerError)
		return
	}

	templateActions, customActions, err := parseJSONActions(device.TemplateActions, device.CustomActions)
	if err != nil {
		log.Printf("error parsing actions: %s", err)
		http.Error(w, "Failed to parse actions", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"template_actions": templateActions,
		"custom_actions":   customActions,
		"id":               device.ID,
	}

	if err := t.Execute(w, data); err != nil {
		log.Printf("error executing feature template: %s", err)
		http.Error(w, "Error executing feature template", http.StatusInternalServerError)
		return
	}
}

func CreateDashboardHandler(w http.ResponseWriter, r *http.Request, db *db.Database) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}
	fmt.Println("Form data:", r.Form)

	dashboardName := r.FormValue("dashboardName")
	if dashboardName == "" {
		http.Error(w, "Dashboard name is required", http.StatusBadRequest)
		return
	}

	r.Form.Del("dashboardName")

	var deviceEntries []model.DeviceInDashboard
	position := 0

	for key, values := range r.Form {
		if strings.HasPrefix(key, "device_action_") {
			deviceIDStr := strings.TrimPrefix(key, "device_action_")
			deviceID, err := strconv.Atoi(deviceIDStr)
			if err != nil {
				http.Error(w, fmt.Sprintf("Invalid device ID %s: %v", deviceIDStr, err), http.StatusBadRequest)
				continue // Skip this iteration
			}

			functionalitiesJSON, err := json.Marshal(values)
			if err != nil {
				http.Error(w, "Error encoding functionalities", http.StatusInternalServerError)
				return
			}

			deviceEntries = append(deviceEntries, model.DeviceInDashboard{
				Device:       model.Device{ID: deviceID},
				ShownActions: string(functionalitiesJSON),
				Position:     position,
			})
			position++
		}
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

	_, err = fmt.Fprintln(w, "Dashboard saved successfully!")
	if err != nil {
		http.Error(w, "Failed to send confirmation: "+err.Error(), http.StatusInternalServerError)
	}
}

func DisplayDashboardHandler(w http.ResponseWriter, r *http.Request, database *db.Database) {
	stringId := r.PathValue("id")
	id, err := strconv.Atoi(stringId)
	if err != nil {
		fmt.Printf("error converting string id to int %s\n", err)
		http.Error(w, "Invalid dashboard ID", http.StatusBadRequest)
		return
	}

	name, devices, err := database.FetchDashboardContents(id)
	if err != nil {
		fmt.Printf("failed to fetch dashboard contents %s\n", err)
		http.Error(w, "Failed to fetch dashboard contents", http.StatusInternalServerError)
		return
	}

	fmt.Println(devices)

	t, err := template.ParseFiles("ui/html/dashboard.gohtml")
	if err != nil {
		fmt.Printf("failed to parse template %s\n", err)
		http.Error(w, "Failed to parse template", http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, map[string]interface{}{
		"Devices": devices,
		"Name":    name,
	})
	if err != nil {
		fmt.Printf("failed to execute template %s\n", err)
		http.Error(w, "Failed to execute template", http.StatusInternalServerError)
		return
	}
}

/*func GetDeviceStateHandler(w http.ResponseWriter, r *http.Request, database *db.Database) {
	stringId := r.PathValue("device_id")
	id, err := strconv.Atoi(stringId)
	if err != nil {
		http.Error(w, "Invalid device ID", http.StatusBadRequest)
		return
	}

	state, err := database.GetDeviceStateByID(id)
	if err != nil {
		http.Error(w, "Failed to fetch device state", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain") // Set content type as text/plain
	_, err = w.Write([]byte(state))
	if err != nil {
		return
	} // Write the state as plain text
}
*/
