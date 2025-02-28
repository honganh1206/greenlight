package main

import (
	"errors"
	"net/http"
	"time"

	"greenlight.honganhpham.net/internal/data"
	"greenlight.honganhpham.net/internal/validator"
)

type registration struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input registration

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	var user *data.User
	// This part is way overkill: Ensure the time taken to send the response is always the same
	err = app.consistentTimeHandler(func() error {
		user = &data.User{
			Name:      input.Name,
			Email:     input.Email,
			Activated: false,
		}

		err = user.Password.Set(input.Password)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return err
		}

		if data.ValidateUser(v, user); !v.Valid() {
			app.failedValidationResponse(w, r, v.Errors)
			return err
		}

		// Test timeout
		// time.Sleep(4 *time.Second)

		return app.models.Users.Insert(user)
	}, minProcessingTime)
	if err != nil {
		switch {
		// FIXME: Change response message to send email even if duplicate
		case errors.Is(err, data.ErrDuplicateEmail):
			app.logger.Error(err, map[string]string{
				"error": "duplicate email registration attempt",
			})
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	token, err := app.models.Token.New(user.ID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Reduce round-trip latency
	// This will be executed CONCURRENTLY
	app.background(func() {
		data := map[string]any{
			"activationToken": token.Plaintext,
			"userID":          user.ID,
		}
		err = app.mailer.Send(user.Email, "user_welcome.tmpl", data)
		if err != nil {
			// Using app.serverResponseError will write the 2nd HTTP response
			// Thus getting a "http: superfluous response.WriteHeader call"
			app.logger.Error(err, nil)
		}
	})

	err = app.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		TokenPlaintext string `json:"token"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}

	v := validator.New()

	if data.ValidateTokenPlaintext(v, input.TokenPlaintext); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Retrieve info of the user associated with the token
	user, err := app.models.Users.GetForToken(data.ScopeActivation, input.TokenPlaintext)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("token", "invalid or expired activation token")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	user.Activated = true

	err = app.models.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Remove associated tokens
	err = app.models.Token.DeleteAllForUser(data.ScopeActivation, user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
