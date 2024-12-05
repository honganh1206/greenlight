package main

// func TestReadIDParam(t *testing.T) {
// 	tl := newTestLogger(t)

// 	app := newTestApplication(t, tl)

// 	tests := []struct {
// 		name          string
// 		urlPath       string
// 		expectedID    int64
// 		expectedError bool
// 	}{
// 		{
// 			name:          "Valid ID",
// 			urlPath:       BaseUrl + MovieV1 + "/123",
// 			expectedID:    123,
// 			expectedError: false,
// 		},
// 		{
// 			name:          "Invalid ID Format",
// 			urlPath:       BaseUrl + MovieV1 + "/abc",
// 			expectedID:    0,
// 			expectedError: true,
// 		},
// 		{
// 			name:          "Negative ID",
// 			urlPath:       BaseUrl + MovieV1 + "/-1",
// 			expectedID:    0,
// 			expectedError: true,
// 		},
// 		{
// 			name:          "Zero ID",
// 			urlPath:       BaseUrl + MovieV1 + "/0",
// 			expectedID:    0,
// 			expectedError: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			r := httptest.NewRequest(http.MethodGet, tt.urlPath, nil)
// 			id, err := app.readIDParam(r)

// 			if tt.expectedError && err == nil {
// 				t.Errorf("expected an error but got none")
// 			}
// 			if !tt.expectedError && err != nil {
// 				assert.NilError(t, err)
// 			}
// 			if id != tt.expectedID {
// 				assert.Equal(t, id, tt.expectedID)
// 			}
// 		})
// 	}
// }

// func TestWriteJSON(t *testing.T) {
// 	tl := newTestLogger(t)

// 	app := newTestApplication(t, tl)

// 	tests := []struct {
// 		name           string
// 		payload        envelope
// 		expectedStatus int
// 		headers        http.Header
// 	}{
// 		{
// 			name: "Valid JSON",
// 			payload: envelope{
// 				"message": "test message",
// 				"data":    map[string]interface{}{"key": "value"},
// 			},
// 			expectedStatus: http.StatusOK,
// 			headers:        http.Header{"X-Custom": []string{"test"}},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			w := httptest.NewRecorder()

// 			err := app.writeJSON(w, tt.expectedStatus, tt.payload, tt.headers)
// 			if err != nil {
// 				assert.NilError(t, err)
// 			}

// 			if w.Code != tt.expectedStatus {
// 				assert.Equal(t, w.Code, tt.expectedStatus)
// 			}

// 			if w.Header().Get("Content-Type") != "application/json" {
// 				t.Error("Content-Type header not set correctly")
// 			}

// 			for k, v := range tt.headers {
// 				if !reflect.DeepEqual(w.Header()[k], v) {
// 					t.Errorf("header %s: got %v; want %v", k, w.Header()[k], v)
// 				}
// 			}
// 		})
// 	}
// }

// func TestReadJSON(t *testing.T) {
// 	tl := newTestLogger(t)

// 	app := newTestApplication(t, tl)

// 	tests := []struct {
// 		name          string
// 		body          string
// 		dst           interface{}
// 		expectedError bool
// 	}{
// 		{
// 			name: "Valid JSON",
// 			body: `{"name": "Test Movie", "year": 2023}`,
// 			dst: &struct {
// 				Name string
// 				Year int
// 			}{},
// 			expectedError: false,
// 		},
// 		{
// 			name:          "Empty Body",
// 			body:          "",
// 			dst:           &struct{}{},
// 			expectedError: true,
// 		},
// 		{
// 			name:          "Invalid JSON Syntax",
// 			body:          `{"name": "Test Movie", "year": }`,
// 			dst:           &struct{}{},
// 			expectedError: true,
// 		},
// 		{
// 			name:          "Multiple JSON Values",
// 			body:          `{"name": "Test"} {"name": "Test2"}`,
// 			dst:           &struct{}{},
// 			expectedError: true,
// 		},
// 		{
// 			name:          "Unknown Field",
// 			body:          `{"unknown_field": "value"}`,
// 			dst:           &struct{}{},
// 			expectedError: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
// 			w := httptest.NewRecorder()

// 			err := app.readJSON(w, r, tt.dst)

// 			if tt.expectedError && err == nil {
// 				t.Errorf("expected an error but got none")
// 			}
// 			if !tt.expectedError && err != nil {
// 				assert.NilError(t, err)
// 			}
// 		})
// 	}
// }
