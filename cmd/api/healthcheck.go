package main

import (
	"net/http"
)

// Implemented as a method for the application struct
// So we do not have to use global vars or closures
// And we can just pass fields of the app struct to the method
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	env := envelope{
		"status":      "available",
		"environment": app.config.env,
		"version":     version,
	}
	err := app.writeJSON(w, http.StatusOK, env, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
