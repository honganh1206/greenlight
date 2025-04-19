# Preflight CORS requests

When a cross-origin request does not meet these conditions, the web browser will trigger a "preflight" request _before the real request_ to determine whether the real cross-origin request is permitted or not

The preflight request always has three components: The `OPTIONS` HTTP method, an `Origin` header, and an `Access-Control-Request-Method`

The `Access-Control-Request-Headers` will not list all headers that the real request use, only headers that are NOT CORS-safe or forbidden

To respond to a preflight request, our `200 OK` response should include:

- A `Access-Control-Allow-Origin` response header
- An `Access-Control-Allow-Methods` header listing HTTP methods which can be used in real cross-origin requests
- An `Access-Control-Allow-Headers` listing headers used in real requests
