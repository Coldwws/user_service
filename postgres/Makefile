include local.migrate.env


LOCAL_BIN := $(abspath $(CURDIR)/../bin)

LOCAL_MIGRATION_DIR=$(MIGRATION_DIR)

LOCAL_MIGRATION_DSN=host=$(PG_HOST) port=$(PG_PORT) user=$(PG_USER) password=$(PG_PASSWORD) dbname=$(PG_DATABASE_NAME) sslmode=$(PG_SSLMODE)


install-deps-goose:
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.14.0

local-migration-status:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATION_DIR} postgres "${LOCAL_MIGRATION_DSN}" status

local-migration-up:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATION_DIR} postgres "${LOCAL_MIGRATION_DSN}" up

local-migration-down:
	${LOCAL_BIN}/goose -dir	${LOCAL_MIGRATION_DIR} postgres "${LOCAL_MIGRATION_DSN}" down