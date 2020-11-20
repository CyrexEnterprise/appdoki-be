OS_NAME := $(shell uname -s | tr A-Z a-z)

build:
	@echo "  >  Building binary..."
	rm -f app-doki-api && go build .

migrate-create:
	./migrations/bin/migrate.$(OS_NAME)-amd64 create \
		-dir ./migrations \
		-ext "sql" \
		-seq \
		"$(MIGRATION_NAME)"

compose:
	source .env && docker-compose up -d --force-recreate --build

test:
	go test ./...

test-verbose:
	go test ./... -v

test-cover:
	go test ./... -cover