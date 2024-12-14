package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Movies MovieModelInterface
}

func NewModels(db *sql.DB) *Models {
	// Return pointer type to ensure we are working with the same instance
	return &Models{
		Movies: MovieModel{DB: db},
	}
}
