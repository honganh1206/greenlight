package mocks

import (
	"greenlight.honganhpham.net/internal/data"
)

func NewMockModels() *data.Models {
	return &data.Models{
		Movies: MockMovieModel{},
	}
}
