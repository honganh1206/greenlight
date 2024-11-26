# Database setup and configurations

## Initialize the database
After composing the PostgreSQL container, we can connect to the PostgreSQL container with this command and check the psql version

```bash
docker exec -it greenlight_postgres bash
psql --version
# Define a new role for root user
su - postgres # Switch the shell user to the postgres system user
psql # Launch the postgres interactive terminal to execute SQL commands
```

As postgres system yser (`postgres=#`) we can check the PostgreSQL user we currently are:

```sql
SELECT current_user;
```

Here we set up the database:

```sql
CREATE DATABASE greenlight;
\c greenlight -- Connect to the greenlight database as user postgres using a meta command
```

We need to create a `greenlight` user without superuser permissions. We will also use an extension of PostgreSQL called `citext` to store user email addresses

```sql
CREATE ROLE greenlight WITH LOGIN PASSWORD 'pa55word';
CREATE EXTENSION IF NOT EXISTS citext;
```

We check the connection inside the container's shell:

```bash
psql --host=localhost --dbname=greenlight --username=greenlight
psql -c 'SHOW config_file;' -- Fine-tune the output file for performance
```

See [this article](https://www.enterprisedb.com/postgres-tutorials/how-tune-postgresql-memory) for tuning PostgreSQL for memory and [this web-based tool](https://pgtune.leopard.in.ua/) so that you don't have to do it manually

## Establish a connection pool
- High level: We configure the data source name (DSN) during runtime so that we can pass the DSN as a command-line flag rather than hard-coding it

## Configuring the Database Connection Pool
- We have two types of connection: in-use and idle. When executing an SQL statement, the connection will be marked as in-use and idle when the task is completed.
- Go will first search for an idle connection to reuse and create a new connection if there is none.
- Go will handle the connections gracefully: Bad connections will be re-tried twice before a new connection is made.
- The `SetMaxIdleConns()` method will default the maximum idle connections to 2. While more idle connections means the app will be more performant, keep in mind that keeping idle connections alive comes with a cost resource-wise.
- Another thing to keep in mind is `MaxIdleConns < MaxOpenConns`
- The `SetConnMaxLifetime()` will by default set NO maximum lifetime - meaning the connections will be reused forever.
- In certain cases, it is useful to *enforce a shorter lifetime* with `ConnMaxLifetime` for example when swapping databases to handle traffic from a load balancer.
- Keep in mind of the frequency the connections will expire: For 100 open connections with each one having a `ConnMaxLifetime` of 1 minute the your app can kill and recreate up to **1.67 connections/s**, which greatly hinders performance.
- The `SetConnMaxIdleTine()` helps us mark connections sitting idle in the connection pool as expired then they will be removed by background cleanup operations.
