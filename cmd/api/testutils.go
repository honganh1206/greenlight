package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"greenlight.honganhpham.net/internal/assert"
	"greenlight.honganhpham.net/internal/logger"
	"greenlight.honganhpham.net/internal/mocks"
)

type testLogger struct {
	*logger.Logger
	Buffer *bytes.Buffer // Capture log output
}

type testServer struct {
	*httptest.Server
}

func newTestLogger(_ *testing.T) *testLogger {
	buffer := &bytes.Buffer{}

	// Configure logger for testing
	cfg := logger.LoggerConfig{
		MinLevel:   logger.LevelInfo,
		StackDepth: 0,     // Disable stack traces for tests
		ShowCaller: false, // Disable caller info for tests
	}

	l := logger.New(buffer, cfg)

	return &testLogger{
		Logger: l,
		Buffer: buffer,
	}
}

func (tl *testLogger) GetLogOutput() string {
	return tl.Buffer.String()
}

func (tl *testLogger) Reset() {
	tl.Buffer.Reset()
}

func newTestApplication(_ *testing.T, tl *testLogger) *application {
	return &application{
		logger: tl.Logger,
		models: mocks.NewMockModels(),
	}
}

func newTestServer(_ *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)
	return &testServer{ts}
}

// Make a GET request to a given URL and return the status code, headers and body
func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, string) {
	rs, err := ts.Client().Get(ts.URL + urlPath)

	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)

	if err != nil {
		t.Fatal(err)
	}

	bytes.TrimSpace(body)

	return rs.StatusCode, rs.Header, string(body)
}

func (ts *testServer) post(t *testing.T, urlPath string, body []byte) (int, http.Header, string) {
	rs, err := ts.Client().Post(ts.URL+urlPath, "application/json", bytes.NewReader(body))

	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()

	respBody, err := io.ReadAll(rs.Body)

	if err != nil {
		t.Fatal(err)
	}

	bytes.TrimSpace(respBody)

	return rs.StatusCode, rs.Header, string(respBody)
}

func (ts *testServer) update(t *testing.T, urlPath string, body []byte) (int, http.Header, string) {
	// Create a new PATCH request (include partial updates)
	req, err := http.NewRequest(http.MethodPatch, ts.URL+urlPath, bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	// Set the content type header
	req.Header.Set("Content-Type", "application/json")

	// Send the request using the test server's client
	rs, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	bytes.TrimSpace(respBody)

	return rs.StatusCode, rs.Header, string(respBody)
}

func (ts *testServer) delete(t *testing.T, urlPath string) (int, http.Header, string) {
	req, err := http.NewRequest(http.MethodDelete, ts.URL+urlPath, nil)
	if err != nil {
		t.Fatal(err)
	}

	rs, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	bytes.TrimSpace(body)
	return rs.StatusCode, rs.Header, string(body)
}

func TestHealthCheck(t *testing.T) {
	tl := newTestLogger(t)

	// Reset the buffer for next test
	t.Cleanup(func() {
		tl.Reset()
	})

	app := newTestApplication(t, tl)

	ts := newTestServer(t, app)

	defer ts.Close()

	code, _, body := ts.get(t, HealthCheckV1)

	assert.Equal(t, code, http.StatusOK)

	assert.Equal(t, body, "OK")
}
