# Managing SQL Query Timeouts

## Mimicking a long-running query

- Go provides *context-aware* variants of `Exec()` and `QueryRow()` as `ExecContext()` and `QueryRowContext()` to deal with cases such as when *a query taking longer to run than expected*. In such cases, we might want to cancel the query -> log an error for further investigation -> return 500.
- We mimic a long-running query by adding `pg_sleep(10` (temporarily) to make the query sleep for 10 seconds before returning the result.
- Tip: We make sure the resources associated with our context are released after the timeout to prevent memory leaks.
- Something super interesting: Our `pq` driver for PostgreSQL is responsible for *sending a cancellation signal to PostgreSQL database* - Our context has a `Done` channel and it will be closed when there is a timeout + `pq` has a gorountine listening to this `Done` channel and it will sed a cancellation signal when the channel `Done` is closed.

## Timeouts outside of PostgreSQL

- Problem: The timeout dealine might be hit **before** the PostgreSQL query even starts.
- Supposed that all of our 25 connections are in use, then any additional queries will be "queued" by `sql.DB` until a connection becomes available. In that case, it is possible that *the timeout deadline will be hit before a connection becomes available*. We can test that by running `go run ./cmd/api -db-max-open-conns=1` and `curl localhost:4000/v1/movies/1 & curl localhost:4000/v1/movies/1 &` on another terminal.
