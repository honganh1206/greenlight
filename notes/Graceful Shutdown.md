# Graceful Shutdown

`Ctrl + C` aka `SIGINT` **immediately** terminates our application, so _in-flight HTTP requests cannot complete_

We will need to:

1. Add shutdown [signals](<https://en.wikipedia.org/wiki/Signal_(IPC)#POSIX_signals>)
2. Use such signals to trigger a graceful shutdown

## Shutdown signals

Some common shutdown signals:

| Signal  | Description                           | Keyboard shortcut | Catchable |
| ------- | ------------------------------------- | ----------------- | --------- |
| SIGINT  | Interrupt from keyboard               | Ctrl + C          | Yes       |
| SIGQUIT | Quit from keyboard                    | Ctrl + \          | Yes       |
| SIGKILL | Immediately kill process              | -                 | No        |
| SIGTERM | Terminate process in _orderly_ manner | -                 | Yes       |

When our app is running, we can use `pgrep -l api` to verify the `api` process exists

We can do `pkill -SIGKILL api` to **immediately** kill the `api` process, or `pkill -SIGTERM api` to terminate it

## What we do?

We spin up a **background goroutine** running for the lifetime of our application. We use `signal.Notify()` to listen to specific signals and relay (push) to a channel for further processing

Note that we use a **buffered** channel here as `Notify()` does NOT wait for a receiver to be available (When the buffer is full either side must block)

Why? Because a signal could be 'missed' if our `quit` channel is **NOT READY** to receive at the exact moment the signal is sent

## What is considered a graceful shutdown?

A graceful shutdown _shuts down the server without interrupting any active connections_. We first close all open listeners -> close all idle connections -> wait for connections to become idle for shut down

What we do: We _instruct our server to stop accepting new HTTP requests_ and _given in-flight request a grace period of 20 seconds to complete_
