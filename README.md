# calc-api

A small HTTP REST service that exposes the arithmetic operations from the
external library **`github.com/strategylippo/calc-lib`**.

> Scaffold only — the service is built by a Builder Epic. See the Epic ticket.

Planned endpoints: `GET /add`, `/subtract`, `/multiply`, `/divide` (all take
`?a=<num>&b=<num>`), plus `GET /health`.
