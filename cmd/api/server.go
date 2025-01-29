package main

import (
	"context"
	"errors"
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

	// Stop accepting new HTTP requests
	// Give in-flight ones 20 seconds to complete
	go func() {
		// Intercepting the signals
		quit := make(chan os.Signal, 1)
		// Listen to upcoming signals and relay (push) to the channel
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		s := <-quit

		app.logger.Info("shutting down server...", map[string]string{
			"signal": s.String(),
		})

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		// Graceful shutdown starts here!
		// Shutdown() will return nil if the graceful shutdown is successful
		// Error if not
		// Whatever the return is, we relay it
		shutdownError <- srv.Shutdown(ctx)

	}()

	app.logger.Info("starting the server", map[string]string{
		"addr": srv.Addr,
		"env":  app.config.env,
	})

	// Calling SHutdown() successfully means ListenAndServe() will IMMEDIATELY return an error
	// It is actually a good thing: That means the graceful shutdown has started!
	err := srv.ListenAndServe()

	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	// If there is an error from the shutdownError channel, it means the graceful shutdown has some errors
	err = <-shutdownError
	if err != nil {
		return err
	}
	app.logger.Info("stopped the server", map[string]string{
		"addr": srv.Addr,
	})

	return nil
}
