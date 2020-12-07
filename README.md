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

These can be run against an existing application or on an isolated Docker environment. 
All necessary commands are available in the Makefile.

It's important that the running application has the environment variable `TEST_MODE` set to `true`.

The tests themselves need two variables: `API_URL` with the URl of where the API is running; `DB_URI` as in the application.

Executing `make integration-tests-compose` will create Docker containers for the API and database and run the tests.

It's also possible to prepare only the containers (`make compose-integration`) and 
leave test running for yourself to, for example, debug the tests in the IDE. 


#### In-Test Authentication

This application uses Google's OpenID Connect Authentication. 
Since Google does not provide test users or similar feature, the application has a test mode 
that skips Bearer token validation but still injects properties in the request context. 
This should be adapted as the tests evolve.

An idea for future improvement is to support an extra OIDC and use it for tests only. 