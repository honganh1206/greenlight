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
		expectedBody   func(*testing.T, []byte)
	}{
		{
			name: "Valid Activation",
			inputJSON: `{
                "token": "P4B3URJZJ2NW5UPZC2OHN4H2NM"
            }`,
			expectedStatus: http.StatusOK,
			expectedBody: func(t *testing.T, body []byte) {
				var response struct {
					User struct {
						ID        int64  `json:"id"`
						Name      string `json:"name"`
						Email     string `json:"email"`
						Activated bool   `json:"activated"`
					} `json:"user"`
				}
				err := json.Unmarshal(body, &response)
				assert.NilError(t, err)
				assert.Equal(t, response.User.Activated, true)
				assert.Equal(t, response.User.Email != "", true)
				assert.Equal(t, response.User.Name != "", true)
				assert.Equal(t, response.User.ID > 0, true)
			},
		},
		{
			name: "Invalid Token Format",
			inputJSON: `{
                "token": "short"
            }`,
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody: func(t *testing.T, body []byte) {
				var response struct {
					Error map[string]string `json:"error"`
				}
				err := json.Unmarshal(body, &response)
				assert.NilError(t, err)
				assert.Equal(t, response.Error["token"] != "", true)
			},
		},
		{
			name: "Missing Token",
			inputJSON: `{
                "token": ""
            }`,
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody: func(t *testing.T, body []byte) {
				var response struct {
					Error map[string]string `json:"error"`
				}
				err := json.Unmarshal(body, &response)
				assert.NilError(t, err)
				assert.Equal(t, response.Error["token"], "must be provided")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.put(t, UserV1+"/activated", []byte(tt.inputJSON))
			assert.Equal(t, code, tt.expectedStatus)

			if tt.expectedBody != nil {
				tt.expectedBody(t, body)
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
