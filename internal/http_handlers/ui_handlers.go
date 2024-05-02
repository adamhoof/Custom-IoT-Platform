package http_handlers

import (
	"html/template"
	"net/http"
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
