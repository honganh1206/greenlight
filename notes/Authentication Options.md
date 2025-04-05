# Authentication Options

We will be comparing between

- [[Basic HTTP authentication]]
- [[Stateful vs Stateless token authentication]]
- [[API key authentication]]
- [[OAuth 2.0 and OpenID Connect]]

## Rule-of-thumb to choose the right authentication option

No 'real' user accounts with slow password hashes -> HTTP basic

No need to store password + users have accounts with 3rd-party providers -> OpenID

Need to delegate authentication -> Stateless authentication tokens

Otherwise -> Stateful authentication tokens

## What we will do

1. Send a JSON request to `POST v1/tokens/authentication` containing the credentials
2. Look up the user record based on the email and check if the password provided is correct
3. Invoke the `New()` method to generate 24-hour token with the scope `authentication`
4. Send the token back to the client in JSON

[[Reading and writing to the request context]]
