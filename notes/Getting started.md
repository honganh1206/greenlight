
### What I did

- Read command-line flags to struct
- Add logger to write messages to standard out stream
- Set up basic Mux (HTTP request multiplexer)
- Add `healthCheckHandler` function as a method for the `application` struct

### APIs

| Method | URL Pattern               | Action                                          |
| ------ | ------------------------- | ----------------------------------------------- |
| GET    | /v1/healthcheck           | Show application health and version information |
| GET    | /v1/movies                | Show the details of all movies                  |
| POST   | /v1/movies                | Create a new movie                              |
| GET    | /v1/movies/:id            | Show the details of a specific movie            |
| PATCH  | /v1/movies/:id            | Update the details of a specific movie          |
| DELETE | /v1/movies/:id            | Delete a specific movie                         |
| POST   | /v1/users                 | Register a new user                             |
| PUT    | /v1/users/activated       | Activate a specific user                        |
| PUT    | /v1/users/password        | Update the password for a specific user         |
| POST   | /v1/tokens/authentication | Generate a new authentication token             |
| POST   | /v1/tokens/password-reset | Generate a new password-reset token             |
| GET    | /debug/vars               | Display application metrics                     |
### Setup

```bash
# Enable module - unique for each project
# This works a a module path as well
go mod init greenlight.honganhpham.net
```

Downloaded dependencies will have its version recorded in the `go.mod` file to ensure reproducible builds. If there is no dependency in the local environment, Go will *automatically download it for you*


```git
$ mkdir -p bin cmd/api internal migrations remote
$ touch Makefile
$ touch cmd/api/main.go
```


```bash
# Folder structure
.
├── bin # Include compiled app binaries (Git ignore this later)
├── cmd
│ └── api # Appplication specific code
│ └── main.go
├── internal #Ancillary packages used by our API (Will be imported in cmd/api)
├── migrations # SQL migration files
├── remote # Configs for prod
├── go.mod
└── Makefile # Automate administrative tasks
```



> [!important] About the `internal` directory
> Any packages under `internal` can *only be imported* by code that is *inside the parent of the `internal` directory*. That is, code in the `internal` can only be used inside the Greenlight project.


Adding handlers as methods of the `application` struct allows us to *make dependencies available without global variables or closures*


```bash
go run ./cmd/api
curl -i localhost:4000/v1/healthcheck
go run ./cmd/api -port=3030 -env=production
```

[[API Versioning]]

## Choosing a router

- I decided to write a custom router following the Regex Table approach - a table of pre-compiled `regexp` objects
