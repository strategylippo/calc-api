# Changelog

All notable changes to calc-api are documented here. The format follows
[Keep a Changelog](https://keepachangelog.com/), and this project uses a
four-part `MAJOR.MINOR.PATCH.MICRO` version scheme.

## [0.1.0.0] - 2026-06-24

### Added
- HTTP REST service exposing arithmetic from `calc-lib`. Call `GET /add`,
  `GET /subtract`, `GET /multiply`, and `GET /divide` with `a` and `b` query
  parameters and get back `{"result": N}`.
- `GET /health` liveness endpoint returning `{"status":"ok"}`.
- Consistent JSON error responses (`{"error": "..."}`) with HTTP 400 for
  missing, non-numeric, or non-finite operands, divide-by-zero, and results
  that aren't a finite number.
- Configurable listen port via the `PORT` environment variable (defaults to
  8080).
- Multi-stage `Dockerfile` (static build on `golang:1.22`, distroless nonroot
  runtime) so the service runs as a small, rootless container.
- README with the endpoint reference, build/test/run instructions, Docker
  usage, and `curl` examples.
- Full-stack integration test that boots a real HTTP server and exercises
  every endpoint plus its error paths, alongside unit tests for handlers,
  routing, and the server entrypoint.

### Security
- The server runs with explicit `ReadHeaderTimeout`, `ReadTimeout`,
  `WriteTimeout`, and `IdleTimeout` so slow-client (Slowloris) connections
  can't tie up server resources indefinitely.

[0.1.0.0]: https://github.com/strategylippo/calc-api/releases/tag/v0.1.0.0
