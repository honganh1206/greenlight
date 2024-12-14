# Managing SQL Query Timeouts

- Go provides *context-aware* variants of `Exec()` and `QueryRow()` as `ExecContext()` and `QueryRowContext()` to deal with cases such as when *a query taking longer to run than expected*. In such cases, we might want to cancel the query -> log an error for further investigation -> return 500.
