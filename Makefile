OS_NAME := $(shell uname -s | tr A-Z a-z)
SHELL := /bin/bash
DB_URI := $(shell grep DB_URI .env | sed 's/DB_URI=//')

### Migrations
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

## Docker Compose
compose-update:
	source .env && docker-compose pull && docker-compose up -d --force-recreate

compose-dev:
	source .env && docker-compose -f docker-compose.dev.yml up -d --force-recreate --build

compose-integration:
	source .env && ENV=test DB_SEED=true docker-compose -f docker-compose.dev.yml up -d --force-recreate --build

### tests
seed:
	go run tests/seed/main.go $(DB_URI)

test:
	go test ./app/...

test-verbose:
	go test ./app/... -v

test-verbose-color:
	go test ./app/... -v | sed ''/PASS/s//$$(printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$$(printf "\033[31mFAIL\033[0m")/''

integration-tests:
	API_URL=http://localhost:4001 go test ./tests/... -v

wait:
	sleep 5

integration-tests-compose: compose-integration wait integration-tests
