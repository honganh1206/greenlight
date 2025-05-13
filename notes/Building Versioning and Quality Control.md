# Building Versioning and Quality Control

```make
# comment
# use tabs, not spaces
target: prerequisite-target-1
    command
    command
```

We can use `make` to pass _named arguments_ when executing a particular rule

```bash
# Pass the name of the migration file as argument
make migration name=create_example_table
```

We can do **namespacing** to create some differentiation between rules and help organize the file like `db/migrations/up` (using `/` is recommended for tab completion)

We can use **prerequisite targets** to ask for confirmation

`make` can be used to _create files on disk_ where the name of a target == name of a file being created. To resolve this we can use `phony target`

```language
.PHONY: target
target: prerequisite-target-1 prerequisite-target-2 ...
command
command
...
```

We can use `.envrc` to replace `.env`

```makefile
# Include variables from the .envrc file
include .envrc
```

[[Quality control code]]

[[Module Proxies and Vendoring]]
