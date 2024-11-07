---
tags:
  - "#study"
  - "#review"
  - "#programming"
  - "#computer"
cssclasses:
  - center-images
---
## What we will do

- Send JSON via REST
- Encode Go native objects into JSON
- Customize how Go objects are encoded into JSON
- Create a reusable helper to send JSON responses to ensure consistency


## Fixed-format JSON

```go
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// As JSON is JUST TEXT
	// We can write the JSON in the same way we write any other text response
	js := `{"status": "available", "environment": %q, "version": %q}`
	js = fmt.Sprintf(js, app.config.env, version)

	w.Header().Set("Content-Type", "application/json")
	// This only accepts byte-typed value so we use casting here
	w.Write([]byte(js))
}

```

## JSON Encoding

- We use `json.Marshal()` to convert Go native objects to JSON as `[]byte`
- Go supports encoding many other native types except `chan, func, complex*`
- We have alternative implementations like `json.Encoder` and `bytes.Buffer` but `json.Marshal()` approach is cleaner
