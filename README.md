# calc-api

A small HTTP REST service that exposes the arithmetic operations from the
external library **`github.com/strategylippo/calc-lib`**.

Every arithmetic endpoint takes two numeric query parameters, `a` and `b`,
delegates the math to calc-lib, and returns JSON. There is also a `/health`
endpoint for liveness checks.

## Endpoints

| Method & path        | Query params      | Success body (200)   |
| -------------------- | ----------------- | -------------------- |
| `GET /health`        | none              | `{"status":"ok"}`    |
| `GET /add`           | `a`, `b` (number) | `{"result":<a+b>}`   |
| `GET /subtract`      | `a`, `b` (number) | `{"result":<a-b>}`   |
| `GET /multiply`      | `a`, `b` (number) | `{"result":<a*b>}`   |
| `GET /divide`        | `a`, `b` (number) | `{"result":<a/b>}`   |

All endpoints are scoped to `GET`. Invalid input returns HTTP `400` with a JSON
error body of the shape `{"error":"<message>"}`. This covers a missing or
non-numeric parameter and division by zero.

## Requirements

- Go 1.22 or newer (see `go.mod`)
- Network access to fetch the `calc-lib` module on first build

## Build, test, run

```sh
# Compile everything (binary lands at ./calc-api, which is gitignored).
go build ./...

# Run the full test suite, including the end-to-end integration test.
go test ./...

# Run the server directly without producing a binary.
go run ./cmd/calc-api
```

### PORT environment variable

The server listens on the port given by the `PORT` environment variable and
falls back to `8080` when `PORT` is unset or empty.

```sh
# Default port 8080.
go run ./cmd/calc-api

# Custom port.
PORT=9090 go run ./cmd/calc-api
```

## Docker

The included multi-stage `Dockerfile` builds a static binary and ships it on a
minimal distroless base. calc-lib is fetched during the build; no library code
is vendored into this repo.

```sh
# Build the image.
docker build -t calc-api .

# Run it, mapping the container port to the host. The server reads PORT.
docker run -e PORT=8080 -p 8080:8080 calc-api
```

Once running, `http://localhost:8080/health` answers from inside the container.

## curl examples

Assuming the server is reachable at `http://localhost:8080`:

```sh
# Health check
curl http://localhost:8080/health
# {"status":"ok"}

# Addition: 2 + 3
curl "http://localhost:8080/add?a=2&b=3"
# {"result":5}

# Subtraction: 5 - 3
curl "http://localhost:8080/subtract?a=5&b=3"
# {"result":2}

# Multiplication: 4 * 3
curl "http://localhost:8080/multiply?a=4&b=3"
# {"result":12}

# Division: 6 / 3
curl "http://localhost:8080/divide?a=6&b=3"
# {"result":2}

# Division by zero -> HTTP 400 with a JSON error body
curl -i "http://localhost:8080/divide?a=1&b=0"
# HTTP/1.1 400 Bad Request
# {"error":"calc: divide by zero"}

# Non-numeric operand -> HTTP 400 with a JSON error body
curl -i "http://localhost:8080/add?a=x&b=3"
# HTTP/1.1 400 Bad Request
# {"error":"query parameter a must be a number"}
```

## Project layout

```
cmd/calc-api        # main: reads PORT, builds the router, serves
internal/server     # NewRouter(): registers /health + arithmetic routes
internal/handlers   # arithmetic handlers calling calc-lib
internal/httpx      # shared JSON response/error helpers
```
