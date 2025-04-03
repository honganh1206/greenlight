# Basic HTTP authentication

We include a `Authorization` header containing the credentials

The credentials needs to be in the format `username:password` and base-64 encoded

We extract the credentials with `Request.BasicAuth()` in Go

It is simple for clients and supported out-of-the-box. We can just send the same header for every request

But it is not a great fit for hashed passwords, since we need to compare the password with the hashed version (slow), and on top of that we need to do so _for every request_
