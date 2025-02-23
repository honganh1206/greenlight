package main

import (
	"net/http"

	"greenlight.honganhpham.net/internal/validator"
)

func (app *application) createActivationTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email string `json:"email"`
	}
	err := app.readJSON(w, r, err)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

}
