# OAuth 2.0 and OpenID Connect

We rely on 3rd-party identity providers like Google or Facebook

> [!WARNING]
> OAuth is NOT an authentication protocol, and we should not use it to authenticate users. To do so, we use OpenID Connect built on top of OAuth 2.0

High level overview of OpenID Connect:

1. When the user authenticates a request, we redirect the user to an 'authentication and consent' form hosted by the identity provider
2. If the user agrees to the form, the identity provider sends a _authorization code_ to our API
3. Our API then sends the authorization code to another endpoint of the identity provider for verification. If valid, the identity provider sends back a JSON response containing an ID token aka JWT
4. From that point we can use the JWT for either stateful or stateless authentication

> [!IMPORTANT]
> The users must have accounts with the identity provider, and they need to complete the authentication and consent step to prove to be human
