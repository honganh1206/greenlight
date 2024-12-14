# Optimistic Concurrency Control

- Changing fields in `input` to pointer type might lead to *race condition* if two clients want to update a movie at a same time. In this case, the `updateMovieHandler` will be running **concurrently** in two different goroutines. And only the latest change will be saved to the database.
- One way to solve this is to the [optimistic locking approach](https://stackoverflow.com/questions/129329/optimistic-vs-pessimistic-locking/129397#129397). In our case, we make sure the 1st update request that reaches our database will succeed, while the second one will return an error message.
- To do so, we make sure that *we filter the record by the version number and update the version number when we make a change*.
- We can add the version number to the `If-Not-Match` or `X-Expected-Version` header to help the client *ensure they are not sending their update request based on outdated information*.
- If you do not want the version identifier to be guess-able, you can use high-entropy random string like UUID for the `version` field.

```sql
UPDATE movies
SET title = $1, year = $2, runtime = $3, genres = $4, version = uuid_generate_v4()
WHERE id = $5 AND version = $6
RETURNING version
```
