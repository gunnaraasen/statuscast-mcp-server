# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
make build      # Build binary to bin/statuscast-mcp-server
make test       # Run all tests (go test ./...)
make lint       # Run golangci-lint
make run-stdio  # Run server with stdio transport
make run-http   # Run server with HTTP transport on :8080

# Run a single test
go test ./internal/config/... -run TestLoad_Defaults
go test ./internal/tools/... -run TestToolNames
```

Required env vars when running:
```bash
STATUSCAST_TOKEN=your-token STATUSCAST_DOMAIN=myco.statuscast.com make run-stdio
```

## Architecture

The server is a thin adapter: it translates MCP tool calls into Statuscast REST API v4 requests.

```
MCP client (Claude Desktop / Kiro)
    ↓ stdio or HTTP (Streamable HTTP)
cmd/server/main.go          — wires config → client → MCP server → transport
    ↓
internal/tools/RegisterAll  — registers 9 tools on the mcp.Server
    ↓
internal/client.Client      — makes authenticated HTTP requests to Statuscast API
```

**MCP SDK pattern** (`github.com/modelcontextprotocol/go-sdk/mcp`): Tools use the generic `mcp.AddTool[In, Out]` with typed input structs. Schema is inferred automatically from `json` and `jsonschema` struct tags — no manual schema writing needed.

```go
type myArgs struct {
    Subject string `json:"subject" jsonschema:"Description shown to LLM"`
}

mcp.AddTool(s, &mcp.Tool{Name: "tool_name", Description: "..."}, func(
    ctx context.Context, _ *mcp.CallToolRequest, args myArgs,
) (*mcp.CallToolResult, any, error) {
    // call client, marshal result, return TextContent
})
```

**Adding a new tool**: add handler + `RegisterXxx` call in the relevant `internal/tools/*.go` file, then call `RegisterXxx` from `tools.go:RegisterAll`. Update the tool count assertion in `tools_test.go`.

**Client** (`internal/client/client.go`): All Statuscast data types and request structs live here alongside the API methods. The `do()` helper handles auth injection, JSON marshal/unmarshal, and maps non-2xx responses to `*APIError`.

**Transport selection**: controlled by `TRANSPORT` env var (`stdio` default, `http` opt-in). The stdio path calls `server.Run(ctx, &mcp.StdioTransport{})` — suitable for Claude Desktop. The HTTP path uses `mcp.NewStreamableHTTPHandler` for multi-session remote deployments.

**In-memory testing**: tools tests use `mcp.NewInMemoryTransports()` + `server.Connect` + `client.Connect` to exercise the full MCP protocol without a real server or network.
