# API key authentication

The idea behind: Every user has a non-expiring secret key associated with their account. The key is hashed similarly to bearer tokens

The key is passed into each request inside the `Authorization` header

After the API receives the key, it can fast hash the key and use it to look up the corresponding user ID

While this is nice for the clients as they do not write code to manage the tokens, the users has another concern: They need to keep their API key secret

Also concerns for developers: We need to implement methods for the users to re-generate their API keys, and generate different keys for different purposes
