package main

import (
	"encoding/json"
	"net/http"
	"testing"

	"greenlight.honganhpham.net/internal/assert"
)

func TestCreateActivationTokenHandler(t *testing.T) {
	tl := newTestLogger(t)

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
		expectedMsg    string
	}{
		{
			name: "Valid Unactivated User",
			inputJSON: `{
                "email": "mock@example.com"
            }`,
			expectedStatus: http.StatusCreated,
			expectedMsg:    "an email will be sent to you containing activation instruction",
		},
		{
			name: "Already Activated User",
			inputJSON: `{
                "email": "activated@example.com"
            }`,
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "Non-existent Email",
			inputJSON: `{
                "email": "nonexistent@example.com"
            }`,
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "Invalid Email Format",
			inputJSON: `{
                "email": "not-an-email"
            }`,
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:           "Missing Email",
			inputJSON:      `{}`,
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "Invalid JSON",
			inputJSON: `{
                "email": "test@example.com"
            `,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Empty Email",
			inputJSON: `{
                "email": ""
            }`,
			expectedStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.post(t, TokenV1+"/activation", []byte(tt.inputJSON))
			assert.Equal(t, code, tt.expectedStatus)

			if tt.expectedStatus == http.StatusCreated {
				var response struct {
					Message string `json:"message"`
				}

				err := json.Unmarshal(body, &response)
				assert.NilError(t, err)
				assert.Equal(t, response.Message, tt.expectedMsg)

				// You might want to add additional checks here:
				// - Verify that an email was attempted to be sent
				// - Check the token was created in the database
				// - Verify the token has the correct scope and expiry
			}
		})
	}
}
