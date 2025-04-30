package middleware

import "net/http"

// Wrapper to see what status code is going to be sent to the client
type CustomResponseWriter struct {
	responseWriter http.ResponseWriter
	StatusCode     int
}

func ExtendedResponseWriter(w http.ResponseWriter) *CustomResponseWriter {
	return &CustomResponseWriter{w, 0}
}

// Implementing the ResponseWriter interface
func (w *CustomResponseWriter) Write(b []byte) (int, error) {
	return w.responseWriter.Write(b)
}

func (w *CustomResponseWriter) Header() http.Header {
	return w.responseWriter.Header()
}

// From the docs: If WriteHeader is not called explicitly
// The first call to Write will trigger an implicit WriteHeader(http.StatusOK).
// Keep that in mind, if w.WriteHeader was not called, the status code should be 200 OK
func (w *CustomResponseWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
	w.responseWriter.WriteHeader(statusCode)
	return
}

func (w *CustomResponseWriter) Done() {
	// if the `w.WriteHeader` wasn't called, set status code to 200 OK
	if w.StatusCode == 0 {
		w.StatusCode = http.StatusOK
	}

	return
}
