OS_NAME := $(shell uname -s | tr A-Z a-z)

compose:
	source .env && docker-compose up -d --force-recreate --build

test:
	go test ./...

test-verbose:
	go test ./... -v

test-cover:
	go test ./... -cover