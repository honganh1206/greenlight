# Cross-Origin Requests

Two URLs share the same origin if they have the same scheme (http/https), host (www/no www) and port

All web browsers implement the **same-origin policy** as a security mechanism

A webpage on one origin can embed certain types of resources from another origin in their HTML

```html
<img src="http://anotherorigin.com/example.png" alt="example image" />
```

A webpage on one origin can send data to a different origin, but it is NOT allowed to receive data from a different origin

The same-origin policy prevents a malicious website on another origin from **reading** information from our website

However, cross-origin sending of data is not prevented by the same policy, so for this we should use `SameSite` cookies and CSRF tokens

Example: Supposed we have `foo.com` running JS, and `bar.com` responding with `data.json`, but _the browser will block the response_ so `foo.com` cannot see it

The `Origin` header is set by the browser to show where the request originates from

The `enableCORS` middleware must be early in the middleware chain, otherwise requests that exceed the rate limit will be blocked by the client's web browser rather than receiving a `429 Too Many Requests`

We can only specify **exactly one origin** in the `Access-Control-Allow-Origin` header

Rule of thumb: If your code makes a decision about what to return based on the content of a request header, you should include that header name in your `Vary` response header - Even of the request did not include that header

> [!IMPORTANT]
> Set `Access-Control-Allow-Credentials: true` in the response if your API endpoints require credentials

[[Preflight CORS requests]]
