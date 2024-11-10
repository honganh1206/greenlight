// Custom Go HTTP router based on a table of regexes

package main

import (
	"context"
	"net/http"
	"regexp"
	"strings"
)

// Empty struct takes zero memory + uniquely identify the key for type safety
type ctxKey struct{}

// Helper to handle path parameters
func getField(r *http.Request, index int) string {
	fields := r.Context().Value(ctxKey{}).([]string)
	return fields[index]
}

type route struct {
	method  string
	regex   *regexp.Regexp
	handler http.HandlerFunc
}

func (app *application) routes() []route {
	return []route{
		newRoute(http.MethodGet, "/v1/healthcheck", app.healthCheckHandler),
		newRoute(http.MethodPost, "/v1/movies", app.createMovieHandler),
		newRoute(http.MethodGet, "/v1/movies/([0-9]+)", app.showMovieHandler),
	}
}

func newRoute(method, pattern string, handler http.HandlerFunc) route {
	return route{method, regexp.MustCompile("^" + pattern + "$"), handler}
}

// Loop through the loop and call the first one that matches both the HTTP method  and the path
func (app *application) serve(w http.ResponseWriter, r *http.Request) {
	var allow []string

	for _, route := range app.routes() {
		matches := route.regex.FindStringSubmatch(r.URL.Path)

		if len(matches) > 0 {
			if r.Method != route.method {
				allow = append(allow, route.method)
				continue
			}

			// Create a new context of HTTP request carrying the URL path params from the matches
			ctx := context.WithValue(r.Context(), ctxKey{}, matches[1:])
			route.handler(w, r.WithContext(ctx))
			return
		}
	}

	// Handle 405
	// Check if any route matches the URL but with an incorrect HTTP method
	if len(allow) > 0 {
		w.Header().Set("Allow", strings.Join(allow, ", "))
		app.methodNotAllowedResponse(w, r)
		return
	}

	// Handle 404
	app.notFoundResponse(w, r)
}
