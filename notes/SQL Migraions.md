# SQL Migrations

## Overview

- When updating our database schema, a migration generates a pair of migration files: One is *up* migration with SQL statements to implement the changes, the other is the *down* migration with SQL statements to reverse/roll-back the changes.
- Each migration pair will be numbered *sequentially* like 001, 002, ... or with Unix timestamp.
- We will use tools/scripts to execute or rollback SQL statements in the migration files.

## Setup go-migrate
- Follow the [online guide](https://github.com/golang-migrate/migrate/blob/master/cmd/migrate/README.md). Note that if you set go in `/usr/local/go/..` then you should move the pre-built binary to `user/local/go..`

## Working with SQL migrations

```bash
# -seq means use sequential numbers
# -ext defines the extension for migration files
# -dir specifies where to store migration files
# a descriptive label
migrate create -seq -ext=.sql -dir=./migrations create_movies_table
```
- Note that working with `NULL` values in Go can be rather awkward, so it is a good practice to set the `NOT NULL` constraints on every table column.
- [This article](https://www.depesz.com/2010/03/02/charx-vs-varcharx-vs-varchar-vs-text/) puts forward some good reasons why we should use `text` instead of `varchar` in certain cases. Some of the reasons might be:
1. Flexibility without length constraints
2. Ease of schema change
3. Simplified usage

- We also need to add some constraints!

```bash
migrate create -seq -ext=.sql -dir=./migrations create_movies_table
```

- After creating migration files, we execute the migration

```bash
# Note that if you are using .env file, you cannot get the environment from the OS
migrate -path=./migrations -database="<your-dsn>" up
```

- Some useful commands after connecting to the DB with psql

```bash
\c <db-name> # connect to the db
\dt # list all databases
\du # list all users
\d <table-name> # see the structure of a table
```
