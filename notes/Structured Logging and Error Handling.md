# Structured Logging and Error Handling

## Structured JSON Log Entries

- For apps which do a lot of logging, we should enforce a consistent structure and format for log entries for easier log entry search & filtering and better integration with 3rd-party analysis and monitoring systems.

- Our custom logger will *write structured log entries in JSON format* - each log entry is a single JSON object.
- The `sync.Mutex` prevents our `Logger` instance from making multiple writes *concurrently*. 

- Note that Go's `http.Server` may also write its own log messages relating to unrecovered panics or problems when accepting/writing to HTTP connections to the **standard logger**, and this will not be added to our custom logger. To solve this, we *pass our custom logger to Go `log.Logger` instance*

### BONUS: My own upgrade for the custom logger

#### Brush up with some terminologies

- A **caller** represents **WHERE** in your code the function was called from. We use `GetCaller()` to get the **immediate caller's info**,  while we use `runtime.CallersFrame()` to get the **context** of the error including program counters, file, line, etc.

- A `frame` is **one level in the stack trace** and the traces show the **sequence** of function calls. Each frame represents one function call in that sequence.

#### Implementation

- Something cool from this custom logger: I added `calldepth` so that we trace the errors better! To do so we have to add the `runtime` library, which has the [runtime.Caller(calldepth)](https://pkg.go.dev/runtime#Caller) method allowing us to *determine how many stack frames to skip when identifying the caller of a function*

As the call stack might look like this:

```js
[0] runtime.Caller()
[1] getCaller()
[2] logger.print()
[3] logger.PrintInfo()  <- This is usually what we want to see
[4] application code    <- Or this, depending on your needs
```

- `skip = 0` would give you info about the `runtime.Caller` line itself
- `skip = 1` would give you info about the `getCaller` function
- `skip = 2` would give you info about the `print` method
- `skip = 3` would give you info about `PrintInfo`
- ...


- Another cool thing: We add a `bufferPool` from the [sync.Pool](https://pkg.go.dev/sync#Pool) library to write our logs and we can reuse it!

- Note that the constant `CALL_DEPTH` explicitly declared and the stack depth from our app configration serve different purposes: The constant serves to *determine how many frames to skip to find the actual caller's location*, while the integer from the configuration is to *determine how many stack frames to collect for the trace*.

## Panic Recovery

- When encountering an error, it is okay to let panics be handled by Go's `http.Server` and then Go will close the underlying HTTP connection + log the error. However, it would be better if *we can send a 500 status response to explain that something has gone wrong*

- In this app we add a middleware to handle panics and also a handler to simulate panic scenarios! with the endpoint `/panic/`
