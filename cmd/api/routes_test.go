package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServe(t *testing.T) {
	tl := newTestLogger(t)

	app := newTestApplication(t, tl)
	// Test cases for different scenarios
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
		// {
		// 	name:           "Matching Route with Incorrect Method",
		// 	method:         "PUT",
		// 	url:            "/test",
		// 	expectedStatus: http.StatusMethodNotAllowed,
		// 	expectedAllow:  "GET, POST",
		// },
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

			// Call the serve method
			app.serve(rec, req)

			// Check status code
			res := rec.Result()
			defer res.Body.Close()
			if res.StatusCode != tc.expectedStatus {
				t.Errorf("expected status %d, got %d", tc.expectedStatus, res.StatusCode)
			}

			// Check Allow header for 405 responses
			if tc.expectedStatus == http.StatusMethodNotAllowed {
				allow := res.Header.Get("Allow")
				if allow != tc.expectedAllow {
					t.Errorf("expected Allow header %q, got %q", tc.expectedAllow, allow)
				}
			}
		})
	}
}
