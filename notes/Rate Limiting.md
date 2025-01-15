# Rate Limiting

We implement rate limiting to _prevent clients from making too many requests_ by creating a middleware.

This middleware will also check how many requests have been received in the last 'N' seconds

What we will do:

1. Implement token-bucket rate-limiter
2. Add a middleware to rate-limit requests (single global limiter -> per-client limiting based on IP)
3. Make limiter behavior **configurable** at runtimne and able to be disabled for testing purposes

## What is a rate limiter BTW?

A **limiter** controls _how frequently events are allowed to happen_

It implements a **token bucket** of size `b` which is _initially full_ and will be refilled at rate `r` tokens per second

## How does our limiter work?

We start with a bucket and `b` tokens in it

Each time we receive a HTTP request, we _remove one token from the bucket_

Every `1/r` seconds, _a token is added back to the bucket_ up to a maximum of `b` total tokens

If the bucket is empty and we receive a HTTP request, we should return a `429` status code response

## IP-based rate limiting

It is better to **separate rate limiter for each client** so that one bad client making too many requests will not impact all the others.

One way to do so is to create an **in-memory** map of rate limiters with IP addresses as keys. Each time a new client makes a request, we _initialize a new rate limiter and add it to the map_

As we will have _multiple goroutines accessing the map concurrently_, we will also need to use `sync.Mutex` to protect access to the map

### Deleting old limiters

Our `clients` map will grow indefinitely and it will take more resources as we add new IP addresses.

We address this by recording the _last seen_ time for each client. We use a background goroutine to check and periodically delete any clients that we haven't been seen recently.

## With distributed applications

It would be better to use a load balancer or reverse proxy with _built-in rate limiting functionality_ like HAProxy or Nginx.

An alternative could be Redis to maintain a request count for clients.
