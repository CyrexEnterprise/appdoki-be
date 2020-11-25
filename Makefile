OS_NAME := $(shell uname -s | tr A-Z a-z)
SHELL := /bin/bash
DB_URI := $(shell grep DB_URI .env | sed 's/DB_URI=//')

compose-update:
	source .env && docker-compose pull && docker-compose up -d --force-recreate

compose-dev:
	source .env && docker-compose -f docker-compose.dev.yml up -d --force-recreate --build

test:
	go test ./...

test-verbose:
	go test ./... -v

test-cover:
	go test ./... -cover


get-migrator:
	mkdir -p migrations/bin
	cd migrations/bin && \
		curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.$(OS_NAME)-amd64.tar.gz | tar xvz

create-migration:
	./migrations/bin/migrate.$(OS_NAME)-amd64 create \
		-dir ./migrations \
		-ext "sql" \
		-seq \
		"$(name)"

migrate-up:
	./migrations/bin/migrate.$(OS_NAME)-amd64 \
		-source file://migrations \
		-database $(DB_URI) up

migrate-down:
	./migrations/bin/migrate.$(OS_NAME)-amd64 \
		-source file://migrations \
		-database $(DB_URI) down