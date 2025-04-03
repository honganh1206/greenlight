package main

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

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

func TestCreateAuthenticationTokenHandler(t *testing.T) {
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
		wantToken      bool
	}{
		{
			name: "Valid Credentials",
			// TODO: No mock password yet
			inputJSON: `{
                "email": "mock@example.com",
                "password": "pa55word"
            }`,
			expectedStatus: http.StatusCreated,
			wantToken:      true,
		},
		{
			name: "Wrong Password",
			inputJSON: `{
                "email": "mock@example.com",
                "password": "wrongpassword"
            }`,
			expectedStatus: http.StatusUnauthorized,
			wantToken:      false,
		},
		{
			name: "Non-existent Email",
			inputJSON: `{
                "email": "nonexistent@example.com",
                "password": "pa55word"
            }`,
			expectedStatus: http.StatusUnprocessableEntity,
			wantToken:      false,
		},
		{
			name: "Invalid Email Format",
			inputJSON: `{
                "email": "not-an-email",
                "password": "pa55word"
            }`,
			expectedStatus: http.StatusUnprocessableEntity,
			wantToken:      false,
		},
		{
			name: "Missing Email",
			inputJSON: `{
                "password": "pa55word"
            }`,
			expectedStatus: http.StatusBadRequest,
			wantToken:      false,
		},
		{
			name: "Missing Password",
			inputJSON: `{
                "email": "mock@example.com"
            }`,
			expectedStatus: http.StatusBadRequest,
			wantToken:      false,
		},
		{
			name:           "Empty JSON",
			inputJSON:      `{}`,
			expectedStatus: http.StatusBadRequest,
			wantToken:      false,
		},
		{
			name: "Invalid JSON",
			inputJSON: `{
                "email": "test@example.com",
                "password": "pa55word"
            `,
			expectedStatus: http.StatusBadRequest,
			wantToken:      false,
		},
		{
			name: "Empty Password",
			inputJSON: `{
                "email": "mock@example.com",
                "password": ""
            }`,
			expectedStatus: http.StatusUnprocessableEntity,
			wantToken:      false,
		},
		{
			name: "Password Too Short",
			inputJSON: `{
                "email": "mock@example.com",
                "password": "short"
            }`,
			expectedStatus: http.StatusUnprocessableEntity,
			wantToken:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, headers, body := ts.post(t, TokenV1+"/authentication", []byte(tt.inputJSON))

			assert.Equal(t, code, tt.expectedStatus)

			if tt.wantToken {
				var response struct {
					Token struct {
						Plaintext string    `json:"plaintext"`
						Expiry    time.Time `json:"expiry"`
					} `json:"authentication_token"`
				}

				err := json.Unmarshal(body, &response)
				assert.NilError(t, err)

				// Check if token is present and valid
				assert.Equal(t, response.Token.Expiry.After(time.Now()), true)

				// Verify content type header
				assert.Equal(t, headers.Get("Content-Type"), "application/json")
			}
		})
	}
}
