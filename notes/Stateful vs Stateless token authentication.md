# Stateful vs Stateless token authentication

AKA bearer token authentication

How it works:

1. The client sends a request containing their credentials
2. The API verifies that the credentials are correct then generate a bearer token, which expires after a period of time
3. The token is included in the `Authorization` header
4. Other APIs check if the token has yet to expire whenver they receive a request

This is better than the basic authentication since we only have to compare the passwords _periodically_

However, this adds complexity for clients since they have to add the logic for caching tokens, monitoring and managing token expiry and refresh token if necessary

## Stateful token

A high-entropy cryptographically-secure random string stored server-side in the database

The big advantage is that the API takes control of the tokens - the downside (not really) is we have to look up the database, but we also need to check the user's activation status and additional information anyway

## Stateless token

Stateless tokens _encode_ the user ID and expiry time **in the token itself** instead of storing the metadata in the database

Most well-known technology is JWT, but there are alternatives like PASETO or Branca

Advantage: The encoding and decoding processes can be done **in memory**, and all the information required to identify the user is within the token already

Downside: Stateless tokens cannot be easily revoked once they are issued. We can partly address that by _changing the secret used for signing the tokens_

Use cases: **Delegated authentication** - The application creating the tokens is **different** from the application consuming them
