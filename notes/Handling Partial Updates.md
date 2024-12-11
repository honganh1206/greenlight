# Handling Partial Updates

- The `updateMovieHandler` should be able to do partial updates of the movie records. For example, we might only want to update the year of a movie and not the genres/runtime/title.
- The key is to differentiate between when a client provides a key-value pair which *has a zero-value value* like `{"title": ""}` which will return a validation error - or when a client does NOT provide a key-value pair in JSON at all which we will 'skip' updating the field instead.
- Note that the default value for pointers would be *nil*, so we can change the fields inside our `input` struct when creating/updating a movie to be of pointer type.
- What if the client explicitly specifies a null value like `curl -X PATCH -d '{"title": null, "year": null}' localhost:4000/v1/movies/4`? In most cases like this, it should be noted in the client documentation that *JSON items with `null` value will be ignored and will remain unchanged*.
