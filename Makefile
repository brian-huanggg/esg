.PHONY: run-backend run-frontend test test-backend  cover-html test-frontend gen gen-backend gen-frontend db-up compose-up compose-down db-exec migrate-up migrate-down lint-frontend format-frontend push

COVER_MIN ?= 80
REGISTRY ?= harbor.pf.internal/lab
# Image tag = git short sha, so a tag always maps back to one commit (never "latest").
TAG ?= $(shell git rev-parse --short HEAD)
PLATFORM ?= linux/amd64
COMPOSE = docker compose --env-file .env -f deploy/compose.yaml
MIGRATE = set -a && . ./.env && set +a && \
	migrate -path backend/migrations \
	-database "mysql://$$DB_USER:$$DB_PASSWORD@tcp(127.0.0.1:3306)/$$DB_NAME"

run-backend:
	cd backend && set -a && . ../.env && set +a && go run ./cmd/api

run-frontend:
	cd frontend && npm run dev

test: test-backend test-frontend

test-backend:
	cd backend && go test -coverprofile=coverage.out ./...
	cd backend && sed -e '/\.gen\.go:/d' -e '/cmd\/api\/main\.go:/d' -e '/shared\/dbtest\//d' coverage.out > coverage.filtered && mv coverage.filtered coverage.out
	cd backend && go tool cover -func=coverage.out | awk -v min=$(COVER_MIN) \
			'/^total:/ { sub(/%/, "", $$3); printf "backend coverage: %s%% (min %s%%)\n", $$3, min; \
			if ($$3 + 0 < min) { print "FAIL: backend coverage below threshold"; exit 1 } }'

cover-html: test-backend
	cd backend && go tool cover -html=coverage.out -o coverage.html

test-frontend:
	cd frontend && npm test

gen: gen-backend gen-frontend

gen-backend:
	cd backend && go generate ./...

gen-frontend:
	cd frontend && npx openapi-typescript ../api/openapi.yaml -o src/api/schema.d.ts

db-up:
	$(COMPOSE) up -d --wait db

compose-up:
	$(COMPOSE) up -d --build --wait

compose-down:
	$(COMPOSE) down

db-exec:
	$(COMPOSE) exec db mariadb -u esg -p

migrate-up:
	$(MIGRATE) up

migrate-down:
	$(MIGRATE) down 1

lint-frontend:
	cd frontend && npm run lint

format-frontend:
	cd frontend && npx prettier --write .

push:
	@test -z "$$(git status --porcelain)" || { echo "FAIL: commit your changes first; the tag $(TAG) must match the pushed code"; exit 1; }
	docker buildx build --platform $(PLATFORM) -t $(REGISTRY)/esg-backend:$(TAG) --push backend
	docker buildx build --platform $(PLATFORM) -t $(REGISTRY)/esg-frontend:$(TAG) --push frontend
	docker buildx build --platform $(PLATFORM) -f backend/Dockerfile.migrate -t $(REGISTRY)/esg-migrate:$(TAG) --push backend
	@echo "pushed tag: $(TAG)"