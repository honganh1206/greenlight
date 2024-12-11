package main

import (
	"net/http"
	"testing"

	"greenlight.honganhpham.net/internal/assert"
)

func TestCreateMovieHandler(t *testing.T) {

	tl := newTestLogger(t)
	app := newTestApplication(t, tl)
	ts := newTestServer(t, app)

	defer ts.Close()

	tests := []struct {
		name           string
		inputJSON      string
		expectedStatus int
	}{
		{
			name: "Valid Input",
			inputJSON: `{
				"title": "Test Movie",
				"year": 2020,
				"runtime": "120 mins",
				"genres": ["drama", "action"]
			}`,
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Missing Title",
			inputJSON: `{
				"year": 2020,
				"runtime": "120 mins",
				"genres": ["drama", "action"]
			}`,
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "Invalid Year",
			inputJSON: `{
				"title": "Test Movie",
				"year": 1800,
				"runtime": "120 mins",
				"genres": ["drama", "action"]
			}`,
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "Too Many Genres",
			inputJSON: `{
				"title": "Test Movie",
				"year": 2020,
				"runtime": "120 mins",
				"genres": ["drama", "action", "comedy", "thriller", "horror", "documentary"]
			}`,
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "Duplicate Genres",
			inputJSON: `{
				"title": "Test Movie",
				"year": 2020,
				"runtime": "120 mins",
				"genres": ["drama", "drama"]
			}`,
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:           "Invalid JSON",
			inputJSON:      `{"title": "Test Movie"`,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, _ := ts.post(t, MovieV1, []byte(tt.inputJSON))
			assert.Equal(t, code, tt.expectedStatus)
		})
	}
}

func TestShowMovieHandler(t *testing.T) {
	tl := newTestLogger(t)
	app := newTestApplication(t, tl)
	ts := newTestServer(t, app)

	defer ts.Close()

	tests := []struct {
		name           string
		urlPath        string
		expectedStatus int
	}{
		{
			name:           "Valid ID",
			urlPath:        "/v1/movies/1",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid ID Format",
			urlPath:        "/v1/movies/abc",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Negative ID",
			urlPath:        "/v1/movies/-1",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, _ := ts.get(t, tt.urlPath)
			assert.Equal(t, code, tt.expectedStatus)
		})
	}
}

func TestUpdateMovieHandler(t *testing.T) {
	tl := newTestLogger(t)
	app := newTestApplication(t, tl)
	ts := newTestServer(t, app)

	defer ts.Close()

	tests := []struct {
		name           string
		urlPath        string
		inputJSON      string
		expectedStatus int
	}{
		{
			name:    "Valid Update",
			urlPath: "/v1/movies/1",
			inputJSON: `{
                "title": "Updated Movie",
                "year": 2021,
                "runtime": "130 mins",
                "genres": ["drama", "sci-fi"]
            }`,
			expectedStatus: http.StatusOK,
		},
		{
			name:    "Invalid ID",
			urlPath: "/v1/movies/abc",
			inputJSON: `{
                "title": "Updated Movie",
                "year": 2021,
                "runtime": "130 mins",
                "genres": ["drama", "sci-fi"]
            }`,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:    "Invalid Input",
			urlPath: "/v1/movies/1",
			inputJSON: `{
                "year": 1800,
                "runtime": "130 mins",
                "genres": ["drama", "drama"]
            }`,
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:           "Malformed JSON",
			urlPath:        "/v1/movies/1",
			inputJSON:      `{"title": "Bad JSON"`,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, _ := ts.update(t, tt.urlPath, []byte(tt.inputJSON))
			assert.Equal(t, code, tt.expectedStatus)
		})
	}
}

func TestDeleteMovieHandler(t *testing.T) {
	tl := newTestLogger(t)
	app := newTestApplication(t, tl)

	ts := newTestServer(t, app)

	defer ts.Close()
	tests := []struct {
		name           string
		urlPath        string
		expectedStatus int
	}{
		{
			name:           "Valid Delete",
			urlPath:        "/v1/movies/1",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid ID Format",
			urlPath:        "/v1/movies/abc",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			code, _, _ := ts.delete(t, tt.urlPath)
			assert.Equal(t, code, tt.expectedStatus)
		})
	}
}
