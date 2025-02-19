package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"greenlight.honganhpham.net/internal/validator"
)

type envelope map[string]any

func (app *application) readIDParam(r *http.Request) (int64, error) {
	params := getField(r, 0)

	// Convert to decimal with a bit size of 64
	id, err := strconv.ParseInt(params, 10, 64)

	if err != nil || id < 1 {
		return 0, errors.New("invalid ID parameter")
	}

	return id, nil
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t") // Add whitespace to encoded JSON

	if err != nil {
		return err
	}

	js = append(js, '\n')

	// No error when ranging over a nil map
	for k, v := range headers {
		w.Header()[k] = v
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains ill-formed JSON")
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)
		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be larger than %d bytes", maxBytesError.Limit)
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}
	}
	// Ensure there is only 1 JSON request body
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must contain only 1 JSON value")
	}
	return nil
}

// Return a string value from a query string
func (app *application) readString(qs url.Values, key, defaultValue string) string {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	return s
}

// Read a string value with commas then split it into a slice of values
func (app *application) readCSV(qs url.Values, key string, defaultValue []string) []string {
	csv := qs.Get(key)

	if csv == "" {
		return defaultValue
	}

	return strings.Split(csv, ",")
}

// Read a string numerical value
func (app *application) readInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(s)

	if err != nil {
		v.AddError(key, "must be an integer value")
		return defaultValue
	}

	return i

}

// Ensure consistent processing time for sensitive operations
func (app *application) consistentTimeHandler(operation func() error, minDuration time.Duration) error {

	startTime := time.Now()

	done := make(chan error, 1)

	go func() {
		done <- operation()
	}()

	// Wait for one of multiple channel operations to complete
	select {
	// Case 1: Operation completes and sends result to done channel
	case err := <-done:
		elapsed := time.Since(startTime)
		// Wait for the remaining time
		if elapsed < minDuration {
			time.Sleep(minDuration - elapsed)
		}
		return err
	// Case 2: Operation takes too long (timeout)
	case <-time.After(minDuration * 2):
		return errors.New("operation timed out")
	}
}

func (app *application) background(fn func()) {
	app.wg.Add(1)
	go func() {
		// Catch panic and log error instead of terminating the application
		defer app.wg.Done()
		defer func() {
			if err := recover(); err != nil {
				app.logger.Error(fmt.Errorf("%s", err), nil)
			}
		}()
	}()
	fn()
}
