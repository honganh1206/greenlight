package main

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"greenlight.honganhpham.net/internal/data"
	"greenlight.honganhpham.net/internal/rate"
	"greenlight.honganhpham.net/internal/validator"
)

func (app *application) recoverPanic(next http.Handler) http.Handler { // Returns a new http.Handler that wraps the anonymous function // http.HandlerFunc converts a function to a Handler
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This will run in the event of a panic as Go unwinds the goroutine stack
		defer func() {
			if err := recover(); err != nil {
				if app.debug {
					stack := debug.Stack()
					// Format the error message based on type
					var errMsg string
					switch v := err.(type) {
					case error:
						errMsg = v.Error()
					default:
						errMsg = fmt.Sprintf("%v", v)
					}

					// Log the full details including stack trace
					app.logger.Error(fmt.Errorf("panic: %s", errMsg), map[string]string{
						"stack_trace": string(stack),
						"url":         r.URL.String(),
						"method":      r.Method,
					})
				}
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (app *application) rateLimit(next http.Handler) http.Handler {

	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex // Only needed for the map operations
		clients = make(map[string]*client)
	)

	// Background goroutine to delete not-seen-recently clients
	go func() {
		for {
			time.Sleep(time.Minute)

			// Lock to prevent check limiter check while the cleaning is taking place
			mu.Lock()

			for ip, client := range clients {
				// Remove clients if they are not seen more than 3 mins
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}

			mu.Unlock()
		}
	}() // Immediate function call operator

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.config.limiter.Enabled {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)

			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}

			// Prevent this code from being executed concurrently
			mu.Lock()

			// Init rate limit for specific
			if _, found := clients[ip]; !found {
				clients[ip] = &client{limiter: rate.New(app.config.limiter)}
			}

			if !clients[ip].limiter.Allow() {
				mu.Unlock()
				app.rateLimitExceedResponse(w, r)
				return
			}

			// IMPORTANT: Unlock the mutex before calling the next handler
			// No defer here - If we do so, the mutex is still locked until the entire handler chain completes aka we complete processing the HTTP request
			mu.Unlock()
		}

		next.ServeHTTP(w, r)

	})
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		authorizationHeader := r.Header.Get("Authorization")

		if authorizationHeader == "" {
			r = app.contextSetUser(r, data.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		token := headerParts[1]

		v := validator.New()

		if data.ValidateTokenPlaintext(v, token); !v.Valid() {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		user, err := app.models.Users.GetForToken(data.ScopeAuthentication, token)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.invalidAuthenticationTokenResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		r = app.contextSetUser(r, user)

		next.ServeHTTP(w, r)
	})
}
