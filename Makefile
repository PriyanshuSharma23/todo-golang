include .envrc

web/dev:
	@go run ./cmd/web -port 4000 -db-dsn '${PSQL_DSN}'

db/psql:
	@psql ${PSQL_DSN}


