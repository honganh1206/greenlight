package main

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"greenlight.honganhpham.net/internal/rate"
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
	cfg := rate.RateLimiterConfig{RequestsPerSecond: 2, BurstSize: 4, QueueSize: 3}
	limiter := rate.New(cfg)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			app.rateLimitExceedResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}
