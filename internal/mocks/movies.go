package mocks

import (
	"time"

	"greenlight.honganhpham.net/internal/data"
)

type MockMovieModel struct{}

var mockMovie = &data.Movie{
	ID:        1,
	CreatedAt: time.Now(),
	Title:     "A sample movie",
	Year:      2000,
	Runtime:   data.Runtime(120),
	Genres:    []string{"drama"},
	Version:   1,
}

func (m MockMovieModel) Insert(movie *data.Movie) error {
	return nil
}

func (m MockMovieModel) Get(id int64) (*data.Movie, error) {
	switch id {
	case 1:
		return mockMovie, nil
	default:
		return nil, data.ErrRecordNotFound
	}
}

func (m MockMovieModel) Update(movie *data.Movie) error {
	return nil
}

func (m MockMovieModel) Delete(id int64) error {
	return nil
}

func (m MockMovieModel) GetAll(
	title string,
	genres []string,
	filters data.Filters) ([]*data.Movie, error) {
	return nil, nil
}
