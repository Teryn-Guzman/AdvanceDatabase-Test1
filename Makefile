## run: run the cmd/api application
include .envrc

## run/api: run the cmd/api application
.PHONY: run/api
run/api:
	@echo  'Running application…'
	@go run ./cmd/api \
	-port=${PORT} \
	-env=${ENVIRONMENT} \
	-limiter-burst=${RATE_LIMITER_BURST} \
	-limiter-rps=${RATE_LIMITER_RPS} \
	-limiter-enabled=${RATE_LIMITER_ENABLED} \
	-db-dsn=${RESTAURANT_DB_DSN} \
	-cors-trusted-origins="${CORS_TRUSTED_ORIGINS}"

## db/psql: connect to the database using psql (terminal)
.PHONY: db/psql
db/psql:
	psql ${RESTAURANT_DB_DSN}

## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}


## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up:
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${RESTAURANT_DB_DSN} up
