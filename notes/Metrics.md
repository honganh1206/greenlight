# Metrics

We use the Go's `expvar` library to collect application metrics during runtime

Some of the returned values

- `TotalAlloc` — Cumulative bytes allocated on the heap (will not decrease).
- `HeapAlloc` — Current number of bytes on the heap.
- `HeapObjects` — Current number of objects on the heap.
- `Sys` — Total bytes of memory obtained from the OS (i.e. total memory reserved by the Go runtime for the heap, stacks, and other internal data structures).
- `NumGC` — Number of completed garbage collector cycles.
- `NextGC` — The target heap size of the next garbage collector cycle (Go aims to keep HeapAlloc ≤ NextGC ).

We can set custom metrics like `expvar.NewString("version").Set(version)`. Note that registering two `expvar` variables with the same name will cause a runtime panic

The `expvar.Publish()` can be used to publish the result of a function

We will have some **request-level metrics** to measure:

- Total number of requests received
- Total number of responses sent
- Total time taken to process all requests in microseconds

[[Visualizing and analyzing metrics]]
