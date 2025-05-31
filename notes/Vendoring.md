# Vendoring

The proxy does NOT guarantee it will store the module forever. Thus we must do `go mod vendor` to **store a complete copy of the source code for 3rd-party packages**

However, this will add bloat to the repository

> There is no easy way to check _if the checksums of the vendored dependencies match those in `go.sum`_. Thus we should _regularly_ run `go mod verify` (check dependencies in module cache match those in `go.sum` ) and `go mod vendor` (copy dependencies to the module cache)

When we deploy our app with Caddy (a reverse-proxy in-front of the app), **alll the requests will come from a single IP address running the Caddy instance**. This will impact our custom rate limiter

Caddy willl have an `X-Forward-For` header for each request containing _the real IP address for the client_

The package `realip` will check the client IP address from either `X-Forward-For` and `X-Real-IP` then fall back to use `r.RemoteAddr`
