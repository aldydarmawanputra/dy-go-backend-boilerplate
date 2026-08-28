.PHONY: help run build tidy test docker-up docker-down migrate-up migrate-down migrate-new

ifneq (,$(wildcard .env))
include .env
export
endif

APP_NAME := app
DB_SSLMODE ?= disable
DATABASE_URL ?= postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)
export DATABASE_URL

help:
	@echo "make run                    - run the API locally"
	@echo "make build                  - build binary to ./bin"
	@echo "make tidy                   - go mod tidy"
	@echo "make test                   - run tests"
	@echo "make docker-up              - start postgres + redis + migrate + app"
	@echo "make docker-down            - stop docker compose (keeps volume)"
	@echo "make migrate-up             - apply dbmate migrations"
	@echo "make migrate-down           - roll back the last migration"
	@echo "make migrate-new name=xxx   - create a new migration file"

run:
	go run ./cmd/api

build:
	go build -o bin/$(APP_NAME) ./cmd/api

tidy:
	go mod tidy

test:
	go test ./...

docker-up:
	docker compose up --build -d

docker-down:
	docker compose down

migrate-up:
	dbmate up

migrate-down:
	dbmate down

migrate-new:
	dbmate new $(name)
