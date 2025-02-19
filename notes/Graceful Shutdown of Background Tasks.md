# Graceful Shutdown of Background Tasks

We can use Go's `sync.WaitGroup` to coordinate goroutines with our graceful shutdown

`sync.WaitGroup` works like a **counter**: When we launch a background goroutine, we increment the counter by 1, and decrement it by 1 when that goroutine finishes.

```go
package main

import (
    "fmt"
    "sync"
)
func main() {
// Declare a new WaitGroup.
    var wg sync.WaitGroup
// Execute a loop 5 times.
    for i := 1; i <= 5; i++ {
// Increment the WaitGroup counter by 1, BEFORE we launch the background routine.
        wg.Add(1)
// Launch the background goroutine.
        go func() {
// Defer a call to wg.Done() to indicate that the background goroutine has
// completed when this function returns. Behind the scenes this decrements
// the WaitGroup counter by 1 and is the same as writing wg.Add(-1).
            defer wg.Done()
        fmt.Println("hello from a goroutine")
        }()
    }
// Wait() blocks until the WaitGroup counter is zero --- essentially blocking until all
// goroutines have completed.
    wg.Wait()
    fmt.Println("all goroutines finished")
}
// Output:
// hello from a goroutine
// hello from a goroutine
// hello from a goroutine
// hello from a goroutine
// hello from a goroutine
// all goroutines finished
```
