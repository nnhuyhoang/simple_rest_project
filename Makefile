include .env
# Configuration
# -------------
PROJECT_ROOT=${shell pwd}
DB_CONTAINER ?= usecase2b_db



.PHONY: stop
stop: ## Stop every service of in the Docker Compose environment
	docker-compose down

.PHONY: init
init: ## Run a PostgreSQL server inside of a Docker Compose environment
	docker-compose up -d
	@echo "Waiting for database connection..."
	@while ! docker exec $(DB_CONTAINER) pg_isready -h ${DB_HOST} -p ${DB_PORT} > /dev/null; do \
		sleep 1; \
	done

.PHONY: migrate-up
migrate-up:
	migrate -path ${PROJECT_ROOT}/migrations/schema -database "postgresql://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" -verbose up

.PHONY: migrate-down
migrate-down:
	migrate -path ${PROJECT_ROOT}/migrations/schema -database "postgresql://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" -verbose down

.PHONY: seed-data
seed-data:
	@docker cp  migrations/seed/seed.sql  ${DB_CONTAINER}:/seed.sql
	@docker exec -t ${DB_CONTAINER} sh -c "PGPASSWORD=${DB_PASS} psql -U ${DB_USER} -d ${DB_NAME} -f /seed.sql"

.PHONY: run
run:
	go run ${PROJECT_ROOT}/cmd/server/main.go

.PHONY: test
test:
	@PROJECT_PATH=$(shell pwd) go test -cover ./...

.PHONY: gen-mock
gen-mock:
	@mockgen -source=./pkg/repo/repo.go -destination=./pkg/repo/mocks/repo.go 

.PHONY: gen-swagger
gen-swagger:
	@swag init -g ./cmd/server/main.go --output ./docs/swagger 
	