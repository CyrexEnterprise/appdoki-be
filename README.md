# appdoki-be

This repository contains the RESTful API server used by our [in-house built company application](https://github.com/Cloudoki/appdoki-rn).

This project, besides having functional expectations, also aims at being a locomotive for technical experiments.

## Technologies

For now the focus is Go and PostgreSQL.

Preferably the most recent versions.

## Development

### API

Aim for API-first development (find the contract in `swaggerui/openapi.yml`).

[API contract](swaggerui/openapi.yml) follows [OpenAPI v3](https://swagger.io/docs/specification/about/).

### Database

Database changes are achieved via migrations.

All migrations have _up_ and _down_ steps, are written in plain SQL and should respect a sequential order.

For convenience, [golang-migrate](https://github.com/golang-migrate/migrate) can also be used to generate the required files.

You can find helper commands in the Makefile for this:
- get-migrator: downloads the migrator bin into `migrations/bin` (this directory is gitignored)
- create-migrations: receives an argument for migration name (ex.:`make create-migration name=alter-users-add-superhero`)
- migrate-up: runs all migrations up
- migrate-down: runs all migrations down

### Setup

- create a PostgreSQL database and user
- create a `.env` file and change accordingly (there is a `.env.sample`)
- run the project (a few options: `go run .`; use your debugger; `make compose-dev`...)

### Integration tests

Integration tests are kept in `./tests` and are developed using Go's testing library and guidelines.

These can be run on an already running application or on an isolated Docker environment. All necessary commands are available in the Makefile.

Prepare for tests by copy `.env` to `.test.env` and change accordingly. Some variables are very important to set while running tests:

```
ENV=test
DB_SEED=true
API_URL=http://localhost:4001
NOTIFIER_TEST_MODE=true
```

Executing `make integration-tests-compose` will create Docker containers for the API and database, seed the database with test data and run the tests.

Seeds are also generated in Go files. Find them in `./seed`.