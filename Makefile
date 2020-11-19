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

compose-up:
	mkdir -p ~/docker/postgres_appdoki
	VOLUME_DIR=~/docker/postgres_appdoki docker-compose up -d

test:
	go test ./...

test-verbose:
	go test ./... -v

test-cover:
	go test ./... -cover