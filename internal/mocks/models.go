package mocks

import "greenlight.honganhpham.net/internal/data"

type MockModels struct {
	Movies data.MovieModelInterface
}

func NewMockModels() *MockModels {
	return &MockModels{
		Movies: MockMovieModel{},
	}
}
