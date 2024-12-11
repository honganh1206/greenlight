# Optimistic Concurrency Control

- Changing fields in `input` to pointer type might lead to *race condition* if two clients want to update a movie at a same time. In this case, the `updateMovieHandler` will be running **concurrently** in two different goroutines. And only the latest change will be saved to the database.
- One way to solve this is to the [optimistic locking approach](https://stackoverflow.com/questions/129329/optimistic-vs-pessimistic-locking/129397#129397)
