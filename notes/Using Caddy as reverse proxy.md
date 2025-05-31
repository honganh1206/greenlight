# Using Caddy as reverse proxy

Caddy uses a `Caddyfile` to listen for HTTP requests from our IP address then acts as a reverse proxy, forwarding the request to port 4000 on our local machine

We can use the `respond` directive to block access to sensitive information e.g., `respond /debug/* "Not Permitted" 403` blocks access to any URL path beginning with `/debug/`

We can still access the production metrics by _opening an SSH tunnel_ between the droplet (port 4000) and our local machine (port 9999 for example)

```bash
ssh -L :9999:45.55.49.87:4000 greenlight@45.55.49.87
```
