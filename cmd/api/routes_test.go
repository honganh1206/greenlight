package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"greenlight.honganhpham.net/internal/assert"
)

func TestServeHTTP(t *testing.T) {

	tl := newTestLogger(t)

	// Reset the buffer for next test
	t.Cleanup(func() {
		tl.Reset()
	})

	app := newTestApplication(t, tl)

	tests := []struct {
		name           string
		method         string
		url            string
		expectedStatus int
		expectedAllow  string
	}{
		{
			name:           "Matching Route with Correct Method",
			method:         "GET",
			url:            "/v1/healthcheck",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Matching Route with Incorrect Method",
			method:         "PUT",
			url:            MovieV1,
			expectedStatus: http.StatusMethodNotAllowed,
			expectedAllow:  "POST, GET",
		},
		{
			name:           "No Matching Route",
			method:         "GET",
			url:            "/notfound",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.url, nil)
			rec := httptest.NewRecorder()

			app.ServeHTTP(rec, req)

			res := rec.Result()
			defer res.Body.Close()
			assert.Equal(t, res.StatusCode, tc.expectedStatus)

			if tc.expectedStatus == http.StatusMethodNotAllowed {
				allow := res.Header.Get("Allow")
				assert.Equal(t, allow, tc.expectedAllow)

			}
		})
	}
}
