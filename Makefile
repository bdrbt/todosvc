all: help
.DEFAULT_GOAL := help

PROJECT_ROOT := $(shell pwd)
TARGET="insofsvc"

# declare local env variables for development
include .env
export

help: ## shows this help
	@cat $(MAKEFILE_LIST) | grep -E '^[a-zA-Z_-]+:.*?## .*$$' | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

format: ## strict formatting with gofumpt
	@echo Adjust formatting...
	@find . -name '*.go' -type f -exec gofumpt -w {} \;

swag-init: ## init swagger
	swag init -g cmd/insofsvc/main.go -o api/docs

lint: format ## Lint project
	golangci-lint run ./...

test-coverage: ## check coverage in default browser
	go test -coverprofile cover.out ./...
	go tool cover -html=cover.out

test: ## runs short tests
	go test -v -count=1 ./...

test-race: ## test with race detector
	go test -race -count=1 ./...

run: ## run development
	@go run cmd/insofsvc/main.go

migration-add: ## migration-add name=$1: create a new database migration
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./database/migrations ${name}

dev-compose-up: ## run docker compose 
	PWD=$(PROJECT_ROOT) docker compose -f deploy/docker/docker-compose.yaml up -d

dev-compose-down: ## remove composed set
	PWD=$(PROJECT_ROOT) docker compose -f deploy/docker/docker-compose.yaml down
