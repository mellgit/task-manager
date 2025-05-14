include .env
LOCAL_BIN:=$(PWD)/bin

# tools
install-deps:
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.23.0
	GOBIN=$(LOCAL_BIN) go install github.com/swaggo/swag/cmd/swag@v1.16.3

clean-cache:
	GOBIN=$(LOCAL_BIN) go clean -cache
	GOBIN=$(LOCAL_BIN) go clean -modcache

# migration build
migration-status:
	$(LOCAL_BIN)/goose -dir $(POSTGRES_MIGRATIONS_PATH) postgres $(POSTGRES_MIGRATIONS_DSN) status -v
migration-add:
	$(LOCAL_BIN)/goose -dir $(POSTGRES_MIGRATIONS_PATH) create $(name) sql
migration-up:
	$(LOCAL_BIN)/goose -dir $(POSTGRES_MIGRATIONS_PATH) postgres $(POSTGRES_MIGRATIONS_DSN) up -v
migration-down:
	$(LOCAL_BIN)/goose -dir $(POSTGRES_MIGRATIONS_PATH) postgres $(POSTGRES_MIGRATIONS_DSN) down -v

# swagger generation
swag:
	$(LOCAL_BIN)/swag init -g cmd/up.go

# docker basic
pa:
	docker ps -a
up:
	docker compose up --build -d
down:
	docker compose down
i:
	docker images
b:
	docker build -t taskmanager .
cleardb:
	rm -r ./postgres_data
r: down cleardb up

lt:
	docker logs -f --tail 100 taskmanager

