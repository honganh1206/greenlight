# Parsing JSON Requests
## JSON Decoding

- We can use either `json.Decoder` or `json.Unmarshal()`
- If the target destination (the pointer as an argument of the `Decode()` method) is a struct, then the struct fields must start with capital letters as they are about to be exported
- Struct tags like `json:"title"` SHOULD be included, as they help Go decodes the values that match the key name. The JSON K-V pairs which cannot be successfully mapped will be ignored.

### Zero values
- If we omit a particular key-value pair in our JSON request body, the value of the left out field will be 0 (supposed the field is of type `int32`). But this begs the question: What if we DELIBERATELY set the value to 0?

## Managing bad requests
- For a public-facing API, the errors must be clear and descriptive enough
- We use `readJSON()` custom helper method to handle errors when decoding JSON values

## Panicking & Returning errors
- The best practice is to *return your errors and handle them gracefully*, but in some cases it is okay to panic and not to stay dogmatic.
- Errors have two classes:
  - The **expected errors** occuring during normal operations are caused by things outside of your program e.g., database query timeout/unavailable network resources/bad user input
  - The **unexpected errors** should not happen during normal operation. In such cases, using `panic()` is widely accepted and Go frequently does this when you make a logical error

## Restricting inputs
- Sometimes we have to deal with *unknown fields* and the best thing we can do is to alert the client to this issue. The `json.Decoder()` can help us with this with `DisallowUnknownFields()` setting.
- As `json.Decoder()` is designed to support *streams* of JSON data, the decoder only reads the 1st JSON value and decodes it while ignoring the second value
- We need to set the upper-limit for the maximum size of the request body, and that can be done with `http.MaxBytesReader()`

## Custom JSON Decoding
- Somestimes we need to *intercept the decoding process* to do some manual conversion. In such cases, Go will call the `UnmarshalJSON()` from the `Unmarshaler` interface to decode JSON values so we might need to implement our own `UnmarshalJSON()`.
