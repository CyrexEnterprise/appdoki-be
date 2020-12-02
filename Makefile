OS_NAME := $(shell uname -s | tr A-Z a-z)
SHELL := /bin/bash
DB_URI := $(shell grep DB_URI .env | sed 's/DB_URI=//')
API_URL := $(shell grep API_URL .test.env | sed 's/API_URL=//')

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
	export $$(cat .env | xargs) && docker-compose pull && docker-compose up -d --force-recreate

compose-dev:
	export $$(cat .env | xargs) && docker-compose -f docker-compose.dev.yml up -d --force-recreate --build

compose-integration:
	export $$(cat .test.env | xargs) && docker-compose -f docker-compose.dev.yml up -d --force-recreate --build

### tests
seed:
	cd seed && go run . $(DB_URI)

test:
	go test ./app/... -v

test-color:
	go test ./app/... -v | sed ''/PASS/s//$$(printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$$(printf "\033[31mFAIL\033[0m")/''

integration-tests:
	API_URL=$(API_URL) go test ./tests/... -v

wait:
	sleep 5

integration-tests-compose: compose-integration wait seed integration-tests

integration-tests-local: seed integration-tests