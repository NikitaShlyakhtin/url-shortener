# Include variables from the .envrc file
include .env

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## run/api: run the cmd/api application
.PHONY: run/api
run/api:
	go run ./cmd/api \
		-base_url=$(BASE_URL) \
		-storage_type=postgres \
		-db-dsn=$(POSTGRES_DSN) \
		-limiter-enabled

## run/mock_client: run the mock client
.PHONY: run/mock_client
run/mock_client:
	go run ./internal/mock_client

## proto/generate: generate the proto files
.PHONY: proto/generate
proto/generate:
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    internal/proto/url_shortener.proto

## flags/api: show the flags for the cmd/api application
.PHONY: flags/api
flags/api:
	go run ./cmd/api -h

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## audit: tidy dependencies, format, vet and test all code
.PHONY: audit
audit:
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify

# ==================================================================================== #
# BUILD
# ==================================================================================== #

## build/api: build the cmd/api application
.PHONY: build/api
build/api:
	@echo 'Building cmd/api...'
	go build -o=./bin/api ./cmd/api
	GOOS=linux GOARCH=amd64 go build -o=./bin/linux_amd64/api ./cmd/api

## build/mock_client: build the internal/mock_client application
.PHONY: build/mock_client
build/mock_client:
	@echo 'Building internal/mock_client...'
	go build -o=./bin/mock_client ./internal/mock_client
	GOOS=linux GOARCH=amd64 go build -o=./bin/linux_amd64/mock_client ./internal/mock_client