SHELL := /bin/bash
BUILD_DATE := `date +%Y%m%d%H%M`
GREEN := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
RESET := $(shell tput -Txterm sgr0)

.PHONY: help

help: ## Show this help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "${YELLOW}%-16s${GREEN}%s${RESET}\n", $$1, $$2}' $(MAKEFILE_LIST)

deps: ## Download dependencies
	@go mod tidy
	@go mod download

deps-devel:
	@brew install jq

tests: ## Run tests
	@go generate ./...
	@go test -cover -race -coverprofile=coverage.txt -covermode=atomic ./...

build: ## Build binary for local operating system
	@go generate ./...
	@go build -ldflags "-s -w -X github.com/debeando/zenit/command.BuildTime=$(BUILD_DATE)" -o zenit main.go
