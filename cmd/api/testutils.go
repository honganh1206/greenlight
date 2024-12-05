package main

import (
	"io"
	"log"
	"testing"

	"greenlight.honganhpham.net/internal/mocks"
)

func newTestLogger(_ *testing.T) *logger {
	return &logger{
		errorLog: log.New(io.Discard, "", 0),
		infoLog:  log.New(io.Discard, "", 0),
	}
}

func newTestApplication(_ *testing.T, tl *logger) *application {
	return &application{
		logger:     tl,
		mockModels: mocks.NewMockModels(),
	}
}
