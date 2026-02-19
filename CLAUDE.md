# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What This Is

An MCP (Model Context Protocol) server for DFIR-IRIS, exposing 88 tools over stdio (JSON-RPC) that let LLM clients interact with DFIR-IRIS incident response cases, alerts, assets, IOCs, timelines, and more.

## Build & Run

```bash
# Build
go build ./cmd/dfir-iris-mcp

# Required env vars
export DFIR_IRIS_URL="https://your-iris-instance.example.com"
export DFIR_IRIS_API_KEY="your-api-key-here"

# Optional: skip TLS verification for dev/demo
export DFIR_IRIS_TLS_SKIP_VERIFY=true

# Manual test via stdio
echo '{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"dfir_iris_system_ping","arguments":{}}}' | \
  DFIR_IRIS_URL=https://your-iris DFIR_IRIS_API_KEY=your-key ./dfir-iris-mcp
```

There are no tests yet. No Makefile — just `go build`.

## Architecture

Entry point: `cmd/dfir-iris-mcp/main.go` — loads config, creates HTTP client, initializes MCP server on stdio, registers all tools.

Three packages under `internal/`:

- **config** — reads `DFIR_IRIS_URL`, `DFIR_IRIS_API_KEY`, and `DFIR_IRIS_TLS_SKIP_VERIFY` from env vars.
- **client** — HTTP client with Bearer token auth. DFIR-IRIS wraps all responses in `{"status","message","data"}` envelopes; the client unwraps these and returns raw `data` as `json.RawMessage`. No typed response structs — raw JSON is passed through for the LLM to interpret.
- **tools** — 14 domain files (`cases.go`, `alerts.go`, `assets.go`, etc.) each registering tools via `mcp.AddTool()`. `register.go` has `RegisterAll()` plus shared helpers: `textResult()`, `errorResult()`, `cidQuery()`, `toBody()`, `toQuery()`.

## Tool Pattern

Every tool follows the same structure:

1. Define an args struct with `json` and `jsonschema` tags
2. Use pointer fields for optional/nullable parameters
3. Register with `mcp.AddTool()` providing name (`dfir_iris_<domain>_<action>`), description, and handler
4. Handler calls `c.Get()` or `c.Post()` and returns `textResult(data)` or `errorResult(err)`

The `toBody()` helper converts arg structs to `map[string]any`, skipping nil pointers and specified exclude fields. The `toQuery()` helper does the same for URL query parameters.

## API Compatibility

Uses DFIR-IRIS legacy API endpoints (not `/api/v2/`) for maximum compatibility across all v2.x versions.

## LimaCharlie Integration

- `playbook_generate_case.py` — LimaCharlie playbook that converts detections into DFIR-IRIS cases (case, alert, asset, timeline event). All secrets (IRIS URL, API key, actions URL) are pulled from LimaCharlie Hive secrets at runtime — never hardcode them.
- `lc-actions-svc/app.py` — Flask microservice for sensor isolation/rejoin from DFIR-IRIS. Requires `LC_API_KEY` and `LC_OID` env vars.
