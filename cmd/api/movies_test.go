package main

// TODO: Uncomment this after setting up mock DB
// func TestCreateMovieHandler(t *testing.T) {

// 	tl := newTestLogger(t)

// 	app := newTestApplication(t, tl)

// 	tests := []struct {
// 		name           string
// 		inputJSON      string
// 		expectedStatus int
// 	}{
// 		{
// 			name: "Valid Input",
// 			inputJSON: `{
// 				"title": "Test Movie",
// 				"year": 2020,
// 				"runtime": "120 mins",
// 				"genres": ["drama", "action"]
// 			}`,
// 			expectedStatus: http.StatusCreated,
// 		},
// 		{
// 			name: "Missing Title",
// 			inputJSON: `{
// 				"year": 2020,
// 				"runtime": "120 mins",
// 				"genres": ["drama", "action"]
// 			}`,
// 			expectedStatus: http.StatusUnprocessableEntity,
// 		},
// 		{
// 			name: "Invalid Year",
// 			inputJSON: `{
// 				"title": "Test Movie",
// 				"year": 1800,
// 				"runtime": "120 mins",
// 				"genres": ["drama", "action"]
// 			}`,
// 			expectedStatus: http.StatusUnprocessableEntity,
// 		},
// 		{
// 			name: "Too Many Genres",
// 			inputJSON: `{
// 				"title": "Test Movie",
// 				"year": 2020,
// 				"runtime": "120 mins",
// 				"genres": ["drama", "action", "comedy", "thriller", "horror", "documentary"]
// 			}`,
// 			expectedStatus: http.StatusUnprocessableEntity,
// 		},
// 		{
// 			name: "Duplicate Genres",
// 			inputJSON: `{
// 				"title": "Test Movie",
// 				"year": 2020,
// 				"runtime": "120 mins",
// 				"genres": ["drama", "drama"]
// 			}`,
// 			expectedStatus: http.StatusUnprocessableEntity,
// 		},
// 		{
// 			name:           "Invalid JSON",
// 			inputJSON:      `{"title": "Test Movie"`,
// 			expectedStatus: http.StatusBadRequest,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			req := httptest.NewRequest(http.MethodPost, MovieV1, bytes.NewBufferString(tt.inputJSON))
// 			req.Header.Set("Content-Type", "application/json")

// 			rr := httptest.NewRecorder()

// 			app.createMovieHandler(rr, req)

// 			assert.Equal(t, rr.Code, tt.expectedStatus)
// 		})
// 	}
// }
