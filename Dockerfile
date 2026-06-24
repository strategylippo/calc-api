# syntax=docker/dockerfile:1

# --- build stage -------------------------------------------------------------
# Pinned to the toolchain declared in go.mod (go 1.22). Building in a full Go
# image keeps the final image free of build tooling.
FROM golang:1.22 AS build
WORKDIR /src

# Resolve modules first so the dependency layer is cached independently of source
# changes. calc-lib is fetched here, not vendored into the repo.
COPY go.mod go.sum ./
RUN go mod download

# Build a static binary so it runs in a minimal scratch-like base with no libc.
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /calc-api ./cmd/calc-api

# --- runtime stage -----------------------------------------------------------
# distroless gives us a tiny image with CA certs and a non-root user, no shell.
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /calc-api /calc-api

# The server reads PORT (default 8080). EXPOSE documents the default; the actual
# port is whatever PORT is set to at runtime.
ENV PORT=8080
EXPOSE 8080

ENTRYPOINT ["/calc-api"]
