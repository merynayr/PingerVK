LOCAL_BIN:=$(CURDIR)/backend/bin


install-golangci-lint:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0

lint:
	$(LOCAL_BIN)/golangci-lint run ./... --config .golangci.pipeline.yaml

docker-build:
	docker compose  --env-file ./backend/.env up --build -d 

docker-run:
	docker compose up -d

include backend/.env

LOCAL_BIN:=$(CURDIR)/backend/bin

LOCAL_MIGRATION_DIR=backend/$(MIGRATION_DIR)
LOCAL_MIGRATION_DSN="host=localhost port=${DB_PORT} dbname=${DB_DB} user=${DB_USER} password=${DB_PASSWORD} sslmode=disable"

local-migration-status:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} status -v

local-migration-up:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} up -v

local-migration-down:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} down -v