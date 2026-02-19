# dfir-iris-mcp

MCP (Model Context Protocol) server for [DFIR-IRIS](https://dfir-iris.org/) — exposing 88 tools that let LLM clients (Claude Desktop, Cursor, Claude Code, etc.) interact with DFIR-IRIS incident response cases, alerts, assets, IOCs, timelines, and more over stdio.

## Quick Start

### Build

```bash
go build ./cmd/dfir-iris-mcp
```

### Configure

Set the required environment variables:

```bash
export DFIR_IRIS_URL="https://your-iris-instance.example.com"
export DFIR_IRIS_API_KEY="your-api-key-here"
```

The API key can be found in DFIR-IRIS under **My Settings > API Key**.

### Run

The server communicates over stdio (JSON-RPC). Add it to your MCP client configuration:

**Claude Desktop / Claude Code** (`claude_desktop_config.json` or MCP settings):

```json
{
  "mcpServers": {
    "dfir-iris": {
      "command": "/path/to/dfir-iris-mcp",
      "env": {
        "DFIR_IRIS_URL": "https://your-iris-instance.example.com",
        "DFIR_IRIS_API_KEY": "your-api-key-here"
      }
    }
  }
}
```

### Test manually

```bash
echo '{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"dfir_iris_system_ping","arguments":{}}}' | \
  DFIR_IRIS_URL=https://your-iris DFIR_IRIS_API_KEY=your-key ./dfir-iris-mcp
```

## Tools (88 total)

| Domain | Tools | Description |
|--------|-------|-------------|
| System | 2 | Ping, version info |
| Settings | 8 | List asset types, IOC types, task statuses, analysis statuses, case states, templates, classifications, evidence types |
| Cases | 9 | List, filter, create, update, delete, close, reopen, summary update, export |
| Alerts | 8 | Filter, get, create, update, delete, escalate, merge, unmerge |
| Assets | 5 | List, get, add, update, delete (case-scoped) |
| Notes | 9 | CRUD for notes and note groups, search (case-scoped) |
| IOCs | 5 | List, get, add, update, delete (case-scoped) |
| Timeline | 5 | List, get, add, update, delete events (case-scoped) |
| Tasks | 5 | List, get, add, update, delete (case-scoped) |
| Evidences | 5 | List, get, add, update, delete (case-scoped) |
| Datastore | 10 | Tree view, file CRUD/move, folder CRUD/move/rename (case-scoped) |
| Comments | 4 | List, add, edit, delete on any case object |
| Users | 5 | List, get, add, update, delete (admin) |
| Groups | 4 | List, add, update, delete (admin) |
| Customers | 4 | List, add, update, delete |

### Naming Convention

All tools follow the pattern `dfir_iris_<domain>_<action>`, e.g.:
- `dfir_iris_cases_list`
- `dfir_iris_alerts_escalate`
- `dfir_iris_assets_add`
- `dfir_iris_timeline_list`

## Architecture

```
cmd/dfir-iris-mcp/main.go         # Entry point
internal/
  config/config.go                 # Env var loading
  client/client.go                 # HTTP client, Bearer auth, envelope unwrap
  tools/
    register.go                    # RegisterAll + helpers
    {domain}.go                    # Tool handlers per domain
```

- **SDK**: Official `github.com/modelcontextprotocol/go-sdk/mcp` (stdio transport)
- **Auth**: Bearer token via `Authorization` header
- **Response handling**: DFIR-IRIS wraps responses in `{"status","message","data"}` — the client unwraps and returns raw `data` JSON
- **No typed response structs**: API responses are passed through as raw JSON text for the LLM to interpret

## Compatibility

Targets the DFIR-IRIS legacy API endpoints which are supported across all DFIR-IRIS v2.x versions. The newer RESTful `/api/v2/` endpoints (IRIS >= v2.4.0) are not used to maximize compatibility.
