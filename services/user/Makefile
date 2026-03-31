include .env
export $(shell sed 's/=.*//' .env)


run:
	@go run cmd/main.go

MIGRATIONS_PATH = ./internal/db/migrate/migrations


compose-up:
	@docker compose up

compose-down:
	@docker compose stop

connect-database:
	@docker exec -it 2fe362e31ed9 psql -U $(DB_USER) -d $(DB_NAME)


.PHONY: migration-create
migration-create:
	@migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(name)

.PHONY: migrate-up
migrate-up:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) up

.PHONY: migrate-down
migrate-down:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) down $(step)

.PHONY: migrate-drop
migrate-drop:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) drop

migrate-version:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) version

migrate-force:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) force $(version)





.PHONY: help
help:
	@echo ""
	@echo "Available commands:"
	@echo ""
	@echo "Application:"
	@echo "  make run                     Run the Go application"
	@echo ""
	@echo "Docker:"
	@echo "  make compose-up              Start docker containers"
	@echo "  make compose-down            Stop docker containers"
	@echo "  make connect-database        Connect to PostgreSQL container"
	@echo ""
	@echo "Migrations:"
	@echo "  make migration-create name=<migration_name>    Create a new migration"
	@echo "  make migrate-up                                Run all pending migrations"
	@echo "  make migrate-down step=<number>                Rollback migrations"
	@echo "  make migrate-drop                              Drop all database tables"
	@echo "  make migrate-version                           Show current migration version"
	@echo "  make migrate-force version=<number>            Force migration version"
	@echo ""
