package main

import (
	"encoding/json"
	"net/http"
	"testing"

	"greenlight.honganhpham.net/internal/assert"
)

func TestRegisterUserHandler(t *testing.T) {
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
		checkTiming    bool
	}{
		// TODO: Sending emails requires more checks
		// {
		// 	name: "Valid Registration",
		// 	inputJSON: `{
		//               "name": "John Doe",
		//               "email": "john@example.com",
		//               "password": "pa55word"
		//           }`,
		// 	expectedStatus: http.StatusCreated,
		// 	expectedMsg:    "an email will be sent to you to complete your registration",
		// 	checkTiming:    true,
		// },
		{
			name: "Duplicate Email",
			inputJSON: `{
                "name": "John Doe",
                "email": "mock@example.com",
                "password": "pa55word"
            }`,
			expectedStatus: http.StatusUnprocessableEntity,
			expectedMsg:    "an email will be sent to you to complete your registration",
			checkTiming:    true,
		},
		{
			name: "Missing Name",
			inputJSON: `{
                "email": "john@example.com",
                "password": "pa55word"
            }`,
			expectedStatus: http.StatusUnprocessableEntity,
			expectedMsg:    "",
			checkTiming:    false,
		}, {
			name: "Missing Password",
			inputJSON: `{
                "name": "John Doe",
                "email": "john@example.com"
            }`,
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "Invalid Email Format",
			inputJSON: `{
                "name": "John Doe",
                "email": "not-an-email",
                "password": "pa55word"
            }`,
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "Password Too Short",
			inputJSON: `{
                "name": "John Doe",
                "email": "john@example.com",
                "password": "short"
            }`,
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "Empty Name",
			inputJSON: `{
                "name": "",
                "email": "john@example.com",
                "password": "pa55word"
            }`,
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			// An email like this already exists. If fails, create one to pass
			name: "Duplicate Email",
			inputJSON: `{
		              "name": "John Doe",
		              "email": "mock@example.com",
		              "password": "pa55word"
		          }`,
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "Invalid JSON",
			inputJSON: `{
                "name": "John Doe",
                "email": "john@example.com"
                "password": "pa55word"
            }`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Empty JSON",
			inputJSON:      `{}`,
			expectedStatus: http.StatusUnprocessableEntity,
		},
		// TODO: Add more validations
		// {
		// 	name: "Name Too Long",
		// 	inputJSON: `{
		//               "name": "This is a very very very very very very very very very very very very long name",
		//               "email": "john@example.com",
		//               "password": "pa55word"
		//           }`,
		// 	expectedStatus: http.StatusUnprocessableEntity,
		// },
		// {
		// 	name: "Email Too Long",
		// 	inputJSON: `{
		//               "name": "John Doe",
		//               "email": "veryveryveryveryveryveryveryveryveryveryverylongemail@veryveryveryveryverylongdomain.com",
		//               "password": "pa55word"
		//           }`,
		// 	expectedStatus: http.StatusUnprocessableEntity,
		// },
	}

	// Store timing measurements for consistent-time operations
	// var timings []time.Duration

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Tracking the total amount of time to make a request
			// start := time.Now()
			code, _, body := ts.post(t, UserV1, []byte(tt.inputJSON))
			// duration := time.Since(start)
			assert.Equal(t, code, tt.expectedStatus)

			if tt.expectedStatus == http.StatusAccepted {
				var response struct {
					Message string `json:"message"`
				}

				err := json.Unmarshal(body, &response)
				assert.NilError(t, err)
				assert.Equal(t, response.Message, tt.expectedMsg)
			}

			// if tt.checkTiming {
			// 	timings = append(timings, duration)
			// }
		})
	}

	// TODO: Since sending emails takes more time, we temporarily comment this out
	// Check timing consistency
	// if len(timings) >= 2 {
	// 	mean := averageDuration(timings)
	// 	for i, timing := range timings {
	// 		variance := timing - mean
	// 		if variance < 0 {
	// 			variance = -variance
	// 		}

	// 		if variance > maxVariance {
	// 			t.Errorf(
	// 				"Timing variance too high for operation %d: got %v, mean %v, variance %v, max allowed %v",
	// 				i, timing, mean, variance, maxVariance,
	// 			)
	// 		}
	// 	}
	// }
}

func TestActivateUserHandler(t *testing.T) {
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
		checkResponse  func(*testing.T, []byte)
	}{
		{
			name: "Valid Activation",
			inputJSON: `{
                "token": "VALIDTOKEN123456789ABCDEFGHIJK"
            }`,
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var response struct {
					User struct {
						Activated bool `json:"activated"`
					} `json:"user"`
				}
				err := json.Unmarshal(body, &response)
				assert.NilError(t, err)
				assert.Equal(t, response.User.Activated, true)
			},
		},
		{
			name: "Invalid Token Format",
			inputJSON: `{
                "token": "tooshort"
            }`,
			expectedStatus: http.StatusUnprocessableEntity,
			checkResponse: func(t *testing.T, body []byte) {
				var response struct {
					Error map[string]string `json:"error"`
				}
				err := json.Unmarshal(body, &response)
				assert.NilError(t, err)
				assert.Equal(t, response.Error["token"] != "", true)
			},
		},
		{
			name: "Expired Token",
			inputJSON: `{
                "token": "EXPIREDTOKEN23456789ABCDEFGHIJK"
            }`,
			expectedStatus: http.StatusUnprocessableEntity,
			checkResponse: func(t *testing.T, body []byte) {
				var response struct {
					Error map[string]string `json:"error"`
				}
				err := json.Unmarshal(body, &response)
				assert.NilError(t, err)
				assert.Equal(t, response.Error["token"], "invalid or expired activation token")
			},
		},
		{
			name: "Missing Token",
			inputJSON: `{
                "token": ""
            }`,
			expectedStatus: http.StatusUnprocessableEntity,
			checkResponse: func(t *testing.T, body []byte) {
				var response struct {
					Error map[string]string `json:"error"`
				}
				err := json.Unmarshal(body, &response)
				assert.NilError(t, err)
				assert.Equal(t, response.Error["token"], "must be provided")
			},
		},
		{
			name:           "Invalid JSON",
			inputJSON:      `{"token": "VALIDTOKEN"`,
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body []byte) {
				var response struct {
					Error map[string]string `json:"error"`
				}
				err := json.Unmarshal(body, &response)
				assert.NilError(t, err)
				assert.Equal(t, response.Error != nil, true)
			},
		},
		{
			name: "Already Activated User",
			inputJSON: `{
                "token": "ACTIVATEDTOKEN3456789ABCDEFGHIJK"
            }`,
			expectedStatus: http.StatusUnprocessableEntity,
			checkResponse: func(t *testing.T, body []byte) {
				var response struct {
					Error map[string]string `json:"error"`
				}
				err := json.Unmarshal(body, &response)
				assert.NilError(t, err)
				assert.Equal(t, response.Error["token"] != "", true)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.post(t, "/v1/users/activate", []byte(tt.inputJSON))
			assert.Equal(t, code, tt.expectedStatus)

			if tt.checkResponse != nil {
				tt.checkResponse(t, body)
			}
		})
	}
}

/*
	HELPER FUNCTIONS
*/

// func measureResponseTime(t *testing.T, ts *testServer, inputJSON string) time.Duration {
// 	start := time.Now()
// 	ts.post(t, UserV1, []byte(inputJSON))
// 	return time.Since(start)
// }

// func averageDuration(durations []time.Duration) time.Duration {
// 	var total time.Duration
// 	for _, d := range durations {
// 		total += d
// 	}
// 	return total / time.Duration(len(durations))
// }
