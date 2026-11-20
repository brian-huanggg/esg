.PHONY: run-backend run-frontend test test-backend test-frontend gen gen-backend gen-frontend

COVER_MIN ?= 80

run-backend:
	cd backend && go run ./cmd/api
run-frontend:
	cd frontend && npm run dev

test: test-backend test-frontend

test-backend:
	cd backend && go test -coverprofile=coverage.out ./...
	cd backend && sed -e '/\.gen\.go:/d' -e '/cmd\/api\/main\.go:/d' coverage.out > coverage.filtered && mv coverage.filtered coverage.out
	cd backend && go tool cover -func=coverage.out | awk -v min=$(COVER_MIN) \
			'/^total:/ { sub(/%/, "", $$3); printf "backend coverage: %s%% (min %s%%)\n", $$3, min; \
			if ($$3 + 0 < min) { print "FAIL: backend coverage below threshold"; exit 1 } }'

test-frontend:
	cd frontend && npm test

gen: gen-backend gen-frontend

gen-backend:
	cd backend && go generate ./...

gen-frontend:
	cd frontend && npx openapi-typescript ../api/openapi.yaml -o src/api/schema.d.ts