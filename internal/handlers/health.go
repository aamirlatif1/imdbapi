package handlers

import (
	"net/http"
)

const version = "1.0.0"

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	data := envelope{
		"status":      "available",
		"environment": "development",
		"version":     version,
	}

	err := writeJSON(w, http.StatusOK, data, nil)

	if err != nil {
		serverErrorResponse(w, r, err)
	}
}
