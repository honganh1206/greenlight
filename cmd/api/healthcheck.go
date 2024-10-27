package main

import (
	"fmt"
	"net/http"
)

// Implemented as a method for the application struct
// So we do not have to use global vars or closures
// And we can just pass fields of the app struct tot he method
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(w, "status: available")
	fmt.Println(w, "environment: %s\n", app.config.env)
	fmt.Println(w, "version: %s\n", version)
}
