# Module Proxies and Vendoring

## Proxies

Go supports _module proxies_ - services that mirror source code from the original repositories

We can check it with `go env` and in the `GOPROXY` variable

The URL `https://proxy.golang.org` points to a _module mirror_ maintained by the Go team at Google, and it contains the source code from many open-source Go packages

We can use the `go get` or `go mod *` command to retrieve the source code. If there is, the mirror will return a zip file, and if not then the mirror will fetch the code from the authoritative repository and forwards it to you

In the worst case scenario (the mirror cannot fetch the code at all), it returns an error response and the `go` tool will fetch the code **directly** from the repository

Benefits:

- Provide protection in case the original repository disappears
- Since we cannot override/remove a package on the mirror once it is stored there, this will prevent the bugs/problems if someone releases an edited version of the same package _with the same version number_
- Fetching modules from the mirror is **much faster** than doing so from repositories

## Vendoring

The proxy does NOT guarantee it will store the module forever. Thus we must do `go mod vendor` to **store a complete copy of the source code for 3rd-party packages**
