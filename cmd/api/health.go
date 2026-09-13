package main

import (
	"net/http"
)

func (app *application) heathHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	data := envelope{
		"status":      "available",
		"environment": app.config.Env,
		"version":     version,
	}

	err := app.writeJSON(w, http.StatusOK, data, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
