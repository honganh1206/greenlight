# Scaling up

Upgrade the droplet vertically (more CPU/memory)

Move PostgreSQL to a separate droplet or use [managed database](https://www.digitalocean.com/products/managed-databases)
If the Go application is the bottleneck -> Profile and optimize your Go application + Run the Go application on multiple droplets

If the PostgreSQL is the bottleneck

- Profile & optimize the database settings
- Cache the results of expensive/frequent database queries
- Move some operations to Redis
- Use read-only database replicas for queries when possible
- Shard the database
