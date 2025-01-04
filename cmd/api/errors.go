package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func (app *application) logError(r *http.Request, err error) {
	app.logger.Error(err, map[string]string{
		"request_method": r.Method,
		"request_url":    r.URL.String(),
	})
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := envelope{"error": message}

	err := app.writeJSON(w, status, env, nil)

	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}
}

func (app *application) debugErrorResponse(w http.ResponseWriter, r *http.Request, status int, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.logError(r, err)
	app.errorResponse(w, r, status, trace)
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	if app.debug {
		app.debugErrorResponse(w, r, http.StatusInternalServerError, err)
		return
	}

	app.logError(r, err)
	message := "the server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	if app.debug {
		app.debugErrorResponse(w, r, http.StatusInternalServerError, err)
		return
	}
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func (app *application) editConflictResponse(w http.ResponseWriter, r *http.Request) {
	msg := "unable to update the record due to an edit conflict, please try again"
	app.errorResponse(w, r, http.StatusConflict, msg)
}
