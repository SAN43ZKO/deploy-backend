.SILENT:

.PHONY: run
run:
	dotenv -f ./.env run -- env ${dev-env-vars} go run cmd/main.go

.PHONY: docs
docs:
	swag init --parseDependency --parseInternal --dir cmd

.PHONY: compose-up
compose-up:
	docker-compose up -d

.PHONY: compose-down
compose-down:
	docker-compose down --remove-orphans
