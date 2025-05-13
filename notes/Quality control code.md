# Quality control code

Some QC commands

```bash
go mod tidy # remove unused dependencies from go.mod and go.sum
go mod verify # check if dependencies have been changed since downloaded
go fmt ./... # format all .go files
go vet ./... # static analysis of code
go test -race -vet=off ./... # run tests with race detector
```
