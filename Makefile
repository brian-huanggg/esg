.PHONY: run-backend run-frontend test cover gen

run-backend:
	cd backend && go run ./cmd/api
run-frontend:
	cd frontend && npm run dev
test:
	cd backend && go test ./...
cover:
	cd backend && go test -cover ./...
gen:
	cd backend && go generate ./...
	cd frontend && npx openapi-typescript ../api/openapi.yaml -o src/api/schema.d.ts