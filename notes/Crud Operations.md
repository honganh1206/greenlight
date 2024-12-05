# CRUD Operations

## Setting up the models

- A good practice is to wrap our models in a parent `Models` struct. This practice allows more readable code in the future like `app.models.Movies.Insert(...)`

## CRUD operations

- We use $N to represent placeholder parameters for the data we want to insert into the query
- Using `RETURNING` clause is PostgreSQL-specific as it allows us to return values from any record manipulated by `INSERT/UPDATE/DELETE` statement
- We can use the same placeholder parameter in multiple positions

```go
// This SQL statement uses the $1 parameter twice, and the value `123` will be used in
// both locations where $1 appears.
stmt := "UPDATE foo SET bar = $1 + $2 WHERE bar = $1"
err := db.Exec(stmt, 123, 456)
if err != nil {
...
}
```

- We can include multiple SQL statements in a single database call, but it must not contain a placeholder parameter. For that, we must use split queries or custom functions in PostgreSQL

```go
stmt := `
UPDATE foo SET bar = true;
UPDATE foo SET baz = false;`
err := db.Exec(stmt)
if err != nil {
...
}
```

- PostgreSQL does not have unsigned integers, so in our Go code we should avoid using using `uint*` types for value we are reading from/writing to PostgreSQL.


### Prepared statements
- While some methods like `Exec()`, `Query()` and `QueryRow()` use prepared statements behind the scene, it is still a good practice to create own own prepared statements for complex SQL statements e.g., those with multiple `JOINS` and got repeated often.
