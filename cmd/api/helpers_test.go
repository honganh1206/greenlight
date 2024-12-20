package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"greenlight.honganhpham.net/internal/assert"
	"greenlight.honganhpham.net/internal/validator"
)

func TestReadIDParam(t *testing.T) {
    tl := newTestLogger(t)
    app := newTestApplication(t, tl)

    tests := []struct {
        name          string
        urlPath       string
        expectedID    int64
        expectedError bool
    }{
        {
            name:          "Valid ID",
            urlPath:       MovieV1 + "/123",
            expectedID:    123,
            expectedError: false,
        },
        {
            name:          "Invalid ID Format",
            urlPath:       MovieV1 + "/abc",
            expectedID:    0,
            expectedError: true,
        },
        {
            name:          "Negative ID",
            urlPath:       MovieV1 + "/-1",
            expectedID:    0,
            expectedError: true,
        },
        {
            name:          "Zero ID",
            urlPath:       MovieV1 + "/0",
            expectedID:    0,
            expectedError: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            r := httptest.NewRequest(http.MethodGet, tt.urlPath, nil)

            // Extract the ID from the URL path
            matches := regexp.MustCompile(`/([0-9a-zA-Z-]+)$`).FindStringSubmatch(tt.urlPath)
            if len(matches) > 1 {
                // Create a new context with the ID parameter
                ctx := context.WithValue(r.Context(), ctxKey{}, []string{matches[1]})
                r = r.WithContext(ctx)
            }

            id, err := app.readIDParam(r)

            if tt.expectedError {
                if err == nil {
                    t.Errorf("expected an error but got none")
                }
            } else {
                assert.NilError(t, err)
                assert.Equal(t, id, tt.expectedID)
            }
        })
    }
}

func TestWriteJSON(t *testing.T) {
	tl := newTestLogger(t)

	app := newTestApplication(t, tl)

	tests := []struct {
		name           string
		payload        envelope
		expectedStatus int
		headers        http.Header
	}{
		{
			name: "Valid JSON",
			payload: envelope{
				"message": "test message",
				"data":    map[string]interface{}{"key": "value"},
			},
			expectedStatus: http.StatusOK,
			headers:        http.Header{"X-Custom": []string{"test"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			err := app.writeJSON(w, tt.expectedStatus, tt.payload, tt.headers)
			if err != nil {
				assert.NilError(t, err)
			}

			if w.Code != tt.expectedStatus {
				assert.Equal(t, w.Code, tt.expectedStatus)
			}

			if w.Header().Get("Content-Type") != "application/json" {
				t.Error("Content-Type header not set correctly")
			}

			for k, v := range tt.headers {
				if !reflect.DeepEqual(w.Header()[k], v) {
					t.Errorf("header %s: got %v; want %v", k, w.Header()[k], v)
				}
			}
		})
	}
}

func TestReadJSON(t *testing.T) {
	tl := newTestLogger(t)

	app := newTestApplication(t, tl)

	tests := []struct {
		name          string
		body          string
		dst           interface{}
		expectedError bool
	}{
		{
			name: "Valid JSON",
			body: `{"name": "Test Movie", "year": 2023}`,
			dst: &struct {
				Name string
				Year int
			}{},
			expectedError: false,
		},
		{
			name:          "Empty Body",
			body:          "",
			dst:           &struct{}{},
			expectedError: true,
		},
		{
			name:          "Invalid JSON Syntax",
			body:          `{"name": "Test Movie", "year": }`,
			dst:           &struct{}{},
			expectedError: true,
		},
		{
			name:          "Multiple JSON Values",
			body:          `{"name": "Test"} {"name": "Test2"}`,
			dst:           &struct{}{},
			expectedError: true,
		},
		{
			name:          "Unknown Field",
			body:          `{"unknown_field": "value"}`,
			dst:           &struct{}{},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			w := httptest.NewRecorder()

			err := app.readJSON(w, r, tt.dst)

			if tt.expectedError && err == nil {
				t.Errorf("expected an error but got none")
			}
			if !tt.expectedError && err != nil {
				assert.NilError(t, err)
			}
		})
	}
}

func TestReadString(t *testing.T) {
    tl := newTestLogger(t)
    app := newTestApplication(t, tl)

    tests := []struct {
        name         string
        queryString  string
        key         string
        defaultValue string
        expected    string
    }{
        {
            name:         "Existing Key",
            queryString:  "name=test",
            key:         "name",
            defaultValue: "default",
            expected:    "test",
        },
        {
            name:         "Missing Key",
            queryString:  "other=value",
            key:         "name",
            defaultValue: "default",
            expected:    "default",
        },
        {
            name:         "Empty Value",
            queryString:  "name=",
            key:         "name",
            defaultValue: "default",
            expected:    "default",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            qs, _ := url.ParseQuery(tt.queryString)
            result := app.readString(qs, tt.key, tt.defaultValue)
            assert.Equal(t, result, tt.expected)
        })
    }
}

func TestReadCSV(t *testing.T) {
    tl := newTestLogger(t)
    app := newTestApplication(t, tl)

    tests := []struct {
        name         string
        queryString  string
        key         string
        defaultValue []string
        expected    []string
    }{
        {
            name:         "Valid CSV",
            queryString:  "genres=action,adventure,comedy",
            key:         "genres",
            defaultValue: []string{"default"},
            expected:    []string{"action", "adventure", "comedy"},
        },
        {
            name:         "Missing Key",
            queryString:  "other=value",
            key:         "genres",
            defaultValue: []string{"default"},
            expected:    []string{"default"},
        },
        {
            name:         "Empty Value",
            queryString:  "genres=",
            key:         "genres",
            defaultValue: []string{"default"},
            expected:    []string{"default"},
        },
        {
            name:         "Single Value",
            queryString:  "genres=action",
            key:         "genres",
            defaultValue: []string{"default"},
            expected:    []string{"action"},
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            qs, _ := url.ParseQuery(tt.queryString)
            result := app.readCSV(qs, tt.key, tt.defaultValue)
            assert.Equal(t, len(result), len(tt.expected))
            for i := range result {
                assert.Equal(t, result[i], tt.expected[i])
            }
        })
    }
}

func TestReadInt(t *testing.T) {
    tl := newTestLogger(t)
    app := newTestApplication(t, tl)

    tests := []struct {
        name         string
        queryString  string
        key         string
        defaultValue int
        expected    int
        expectError bool
    }{
        {
            name:         "Valid Integer",
            queryString:  "page=5",
            key:         "page",
            defaultValue: 1,
            expected:    5,
            expectError: false,
        },
        {
            name:         "Missing Key",
            queryString:  "other=value",
            key:         "page",
            defaultValue: 1,
            expected:    1,
            expectError: false,
        },
        {
            name:         "Invalid Integer",
            queryString:  "page=abc",
            key:         "page",
            defaultValue: 1,
            expected:    1,
            expectError: true,
        },
        {
            name:         "Empty Value",
            queryString:  "page=",
            key:         "page",
            defaultValue: 1,
            expected:    1,
            expectError: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            qs, _ := url.ParseQuery(tt.queryString)
            v := validator.New()
            result := app.readInt(qs, tt.key, tt.defaultValue, v)
            assert.Equal(t, result, tt.expected)

            if tt.expectError {
                if len(v.Errors) == 0 {
                    t.Error("expected validation error but got none")
                }
            } else {
                if len(v.Errors) > 0 {
                    t.Error("unexpected validation error")
                }
            }
        })
    }
}
