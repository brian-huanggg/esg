.PHONY: run-backend run-frontend test test-backend test-frontend gen

COVER_MIN ?= 80

run-backend:
	cd backend && go run ./cmd/api
run-frontend:
	cd frontend && npm run dev

test: test-backend test-frontend

test-backend:
	cd backend && go test -coverprofile=coverage.out ./...
	# Excluded from the threshold:
	# - *.gen.go: generated code is not ours to test.
	# - cmd/api/main.go: process wiring (ListenAndServe, os.Exit); main_test.go
	#   still runs, it just doesn't count. Keep logic out of main for this reason.
	cd backend && sed -i -e '/\.gen\.go:/d' -e '/cmd\/api\/main\.go:/d' coverage.out
	cd backend && go tool cover -func=coverage.out | awk -v min=$(COVER_MIN) \
			'/^total:/ { sub(/%/, "", $$3); printf "backend coverage: %s%% (min %s%%)\n", $$3, min; \
			if ($$3 + 0 < min) { print "FAIL: backend coverage below threshold"; exit 1 } }'

test-frontend:
	cd frontend && npm test

gen:
	cd backend && go generate ./...
	cd frontend && npx openapi-typescript ../api/openapi.yaml -o src/api/schema.d.ts