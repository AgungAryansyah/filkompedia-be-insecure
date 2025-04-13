include .env
export

DB_URL := postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

MIGRATE_DIR := sql/migrations

createmigration:
	migrate create -ext=sql -dir=$(MIGRATE_DIR) -seq init

migrate:
	migrate -path=$(MIGRATE_DIR) -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path=$(MIGRATE_DIR) -database "$(DB_URL)" -verbose down

.PHONY: migrate migratedown createmigration
