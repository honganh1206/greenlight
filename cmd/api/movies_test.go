package main

import (
	"net/http"
	"testing"
	"time"

	"greenlight.honganhpham.net/internal/assert"
	"greenlight.honganhpham.net/internal/data"
)

func TestCreateMovieHandler(t *testing.T) {

	tl := newTestLogger(t)

	// Reset the buffer for next test
	t.Cleanup(func() {
		tl.Reset()
	})

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
		{
			name: "Empty Runtime",
			inputJSON: `{
            "title": "Test Movie",
            "year": 2020,
            "runtime": "",
            "genres": ["drama"]
        }`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Future Year",
			inputJSON: `{
            "title": "Test Movie",
            "year": 2525,
            "runtime": "120 mins",
            "genres": ["drama"]
        }`,
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "Empty Genres Array",
			inputJSON: `{
            "title": "Test Movie",
            "year": 2020,
            "runtime": "120 mins",
            "genres": []
        }`,
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "Invalid Runtime Format",
			inputJSON: `{
            "title": "Test Movie",
            "year": 2020,
            "runtime": "invalid",
            "genres": ["drama"]
        }`,
			expectedStatus: http.StatusBadRequest,
		},
	}

	token, err := app.models.Tokens.New(1, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, _ := ts.post(t, MovieV1, token.Plaintext, []byte(tt.inputJSON))
			assert.Equal(t, code, tt.expectedStatus)
		})
	}
}

func TestShowMovieHandler(t *testing.T) {
	tl := newTestLogger(t)

	// Reset the buffer for next test
	t.Cleanup(func() {
		tl.Reset()
	})

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
		{
			name:           "Zero ID",
			urlPath:        "/v1/movies/0",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Very Large ID",
			urlPath:        "/v1/movies/999999999",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Invalid Path",
			urlPath:        "/v1/movies/1/invalid",
			expectedStatus: http.StatusNotFound,
		},
	}

	token, err := app.models.Tokens.New(1, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, _ := ts.get(t, tt.urlPath, token.Plaintext)
			assert.Equal(t, code, tt.expectedStatus)
		})
	}
}

func TestUpdateMovieHandler(t *testing.T) {
	tl := newTestLogger(t)

	// Reset the buffer for next test
	t.Cleanup(func() {
		tl.Reset()
	})

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
		{
			name:           "Empty Request Body",
			urlPath:        "/v1/movies/1",
			inputJSON:      `{}`,
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:    "Invalid Runtime Format",
			urlPath: "/v1/movies/1",
			inputJSON: `{
            "title": "Updated Movie",
            "year": 2021,
            "runtime": "invalid",
            "genres": ["drama"]
        }`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:    "Future Year",
			urlPath: "/v1/movies/1",
			inputJSON: `{
            "title": "Updated Movie",
            "year": 2525,
            "runtime": "120 mins",
            "genres": ["drama"]
        }`,
			expectedStatus: http.StatusUnprocessableEntity,
		},
	}

	token, err := app.models.Tokens.New(1, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, _ := ts.patch(t, tt.urlPath, token.Plaintext, []byte(tt.inputJSON))
			assert.Equal(t, code, tt.expectedStatus)
		})
	}
}

func TestDeleteMovieHandler(t *testing.T) {
	tl := newTestLogger(t)

	// Reset the buffer for next test
	t.Cleanup(func() {
		tl.Reset()
	})

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
		{
			name:           "Zero ID",
			urlPath:        "/v1/movies/0",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Invalid Path",
			urlPath:        "/v1/movies/1/invalid",
			expectedStatus: http.StatusNotFound,
		},
	}

	token, err := app.models.Tokens.New(1, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, _ := ts.delete(t, tt.urlPath, token.Plaintext)
			assert.Equal(t, code, tt.expectedStatus)
		})
	}
}
