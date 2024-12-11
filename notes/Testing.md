# Testing

## Mocking
- We create a `NewMockModels()` function to return a `*Models` type. As our `Movies` field inside the `Models` struct is of type `MovieModelInterface` interface, we must ensure that both of our `MovieModel` and `MockMovieModel` structs implements the `MovieModelInterface`

## E2E Testing
- As HTTP handlers are not actually used in isolation most of the times, we need something to run E2E tests on which encompass our routing, middleware and handlers. FOr this reason, we can use `httptest.Server` as a memory HTTP server.
