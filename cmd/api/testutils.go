package main

import (
	"io"
	"log"
	"testing"
)

func newTestLogger(t *testing.T) *logger {
	return &logger{
		errorLog: log.New(io.Discard, "", 0),
		infoLog:  log.New(io.Discard, "", 0),
	}
}

func newTestApplication(t *testing.T, tl *logger) *application {

	return &application{
		logger: tl,
	}
}
