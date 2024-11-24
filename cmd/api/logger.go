package main

import (
	"log"
	"os"
)

type logger struct {
	infoLog  *log.Logger
	errorLog *log.Logger
	fileLog  *log.Logger
}

func newLogger(f *os.File) *logger {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	fileLog := log.New(f, "INFO\t", log.Ldate|log.Ltime)

	return &logger{
		infoLog:  infoLog,
		errorLog: errorLog,
		fileLog:  fileLog,
	}
}
