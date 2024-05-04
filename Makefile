.PHONY: generate-sql-models
.PHONY: install-tools

include .env

install-tools:
	go install github.com/go-jet/jet/v2/cmd/jet@latest
	go install github.com/dmarkham/enumer

generate-sql-models:
	jet -dsn=postgresql://$(DATABASE_USER):$(DATABASE_PASSWORD)@$(DATABASE_HOST):$(DATABASE_PORT)/$(DATABASE_NAME)?sslmode=$(DATABASE_SSL) -path=internal/repository/pgsql/.gen