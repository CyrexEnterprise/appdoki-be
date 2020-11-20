OS_NAME := $(shell uname -s | tr A-Z a-z)

SHELL := /bin/bash

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