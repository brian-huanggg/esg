# --- Build stage: full Go toolchain, discarded after compiling ---
FROM golang:1.26 AS build
WORKDIR /src

# Copy dependency manifests first so this layer is cached until go.mod/go.sum change.
COPY go.mod go.sum ./
RUN go mod download

COPY . .
# CGO_ENABLED=0 produces a static binary; distroless/static has no libc.
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api

# --- Runtime stage: no shell, no package manager, runs as non-root ---
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/api /api
EXPOSE 8080
ENTRYPOINT ["/api"]
