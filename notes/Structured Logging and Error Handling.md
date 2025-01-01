# Structured Logging and Error Handling

## Structured JSON Log Entries

- For apps which do a lot of logging, we should enforce a consistent structure and format for log entries for easier log entry search & filtering and better integration with 3rd-party analysis and monitoring systems.

- Our custom logger will *write structured log entries in JSON format* - each log entry is a single JSON object.

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

- Here is an interesting (and might be unnecessary) thing: We declare a `bufferPool` from the [sync.Pool](https://pkg.go.dev/sync#Pool) library

- Note that the constant `CALL_DEPTH` explicitly declared and the stack depth from our app configration serve different purposes: The constant serves to *determine how many frames to skip to find the actual caller's location*, while the integer from the configuration is to *determine how many stack frames to collect for the trace*
