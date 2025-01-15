package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port), // String formatting
		Handler:      app.recoverPanic(app.rateLimit(http.HandlerFunc(app.ServeHTTP))),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second, // TODO: Hardcoded values here
		WriteTimeout: 30 * time.Second,
		ErrorLog:     log.New(app.logger, "", 0),
	}

	shutdownError := make(chan error)

	go func() {
		// Intercepting the signals
		quit := make(chan os.Signal, 1)
		// Listen to upcoming signals and relay (push) to the channel
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		s := <-quit

		app.logger.Info("shutting down server...", map[string]string{
			"signal": s.String(),
		})
		// Graceful shutdown here!
		os.Exit(0)

	}()

	app.logger.Info("Starting the server", map[string]string{
		"addr": srv.Addr,
		"env":  app.config.env,
	})

	return srv.ListenAndServe()
}
