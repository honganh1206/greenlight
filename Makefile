include .envrc
# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	# Ask then validate the response to true/false
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

.PHONY: build/api
build/api:
	@echo 'Building cmd/api...'
	# Reduce binary size by stripping the DWARF debugging information
	go build -ldflags='-s' -o=./bin/api ./cmd/api
	# Cross-compilation
	GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=./bin/linux_amd64/api ./cmd/api

## run/api: run the cmd/api application
.PHONY: run/api
run/api:
	go run ./cmd/api -dsn=${GREENLIGHT_DB_DSN}

## db/psql: connect to the database using psql
.PHONY: db/psql
db/psql:
	psql ${GREENLIGHT_DB_DSN}

## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration file for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## db/migrations/up: apply all migrations
.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${GREENLIGHT_DB_DSN} up

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #
## audit: tidy dependencies and format, vet and test all code
.PHONY: audit
audit:
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...

.PHONY: vendor
vendor:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Vendoring dependencies...'
	go mod vendor

# ==================================================================================== #
# PRODUCTION
# ==================================================================================== #

## production/connect: connect to the production server
.PHONY: production/connect
production/connect:
	ssh greenlight@${GREENLIGHT_PRODUCTION_HOST_IP}

## production/deploy/api: deploy the api to production
.PHONY: production/deploy/api
production/deploy/api:
	rsync -P ./bin/linux_amd64/api greenlight@${GREENLIGHT_PRODUCTION_HOST_IP}:~
	rsync -rP --delete ./migrations greenlight@${GREENLIGHT_PRODUCTION_HOST_IP}:~
	rsync -P ./remote/production/api.service greenlight@${GREENLIGHT_PRODUCTION_HOST_IP}:~
	rsync -P ./remote/production/Caddyfile greenlight@${GREENLIGHT_PRODUCTION_HOST_IP}:~
	# Actual command running on the droplet: 'migrate -path ~/migrations -database $GREENLIGHT_DB_DSN up'
	ssh -t greenlight@${GREENLIGHT_PRODUCTION_HOST_IP} '\
	    migrate -path ~/migrations -database $$GREENLIGHT_DB_DSN up \
		&& sudo mv ~/api.service /etc/systemd/system/ \
        && sudo systemctl enable api \
        && sudo systemctl restart api \
        && sudo mv ~/Caddyfile /etc/caddy/ \
        && sudo systemctl reload caddy \
'
