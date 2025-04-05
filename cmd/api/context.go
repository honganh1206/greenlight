package main

import (
	"context"
	"net/http"

	"greenlight.honganhpham.net/internal/data"
)

type contextKey string

// To be used as a key to get and set user information in the request context
const userContextKey = contextKey("user")

// Return a copy of the context with user data
func (app *application) contextSetUser(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func (app *application) contextGetUser(r *http.Request) *data.User {
	user, ok := r.Context().Value(userContextKey).(*data.User)
	if !ok {
		panic("missing user value in request context")
	}

	return user
}
