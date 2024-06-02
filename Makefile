SHELL := /bin/bash
include .env
export
export APP_NAME := $(basename $(notdir $(shell pwd)))

.PHONY: help
help: ## display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: build
build: ## go build
	@go build -v ./... && go clean

.PHONY: fmt
fmt: ## go format
	@go fmt ./...

.PHONY: mod
mod: ## go mod tidy & go mod vendor
	@go mod tidy

.PHONY: update
update: ## go modules update
	@go get -u -t ./...
	@go mod tidy

.PHONY: gen
gen: ## Generate code.
	@oapi-codegen -generate std-http -package api api/openapi.yaml > pkg/api/server.gen.go
	@oapi-codegen -generate types -package api api/openapi.yaml > pkg/api/types.gen.go
	@go generate ./...
	@go mod tidy
