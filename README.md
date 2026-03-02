# statuscast-mcp-server

A Model Context Protocol (MCP) server that exposes [Statuscast](https://statuscast.com) status page management as LLM-callable tools. Use it with Claude Desktop, Kiro, or any MCP-compatible client to create incidents, inspect component status, manage subscribers, and work with content templates — all via the Statuscast REST API v4.

## Features

- **9 tools** across 4 API domains:
  - **Incidents**: create, get, update, search
  - **Components**: list all, get historical uptime
  - **Subscribers**: create
  - **Content Templates**: list, create
- **Two transports**: stdio (default, for local use) and Streamable HTTP (for remote deployments)
- **Single binary**: no runtime dependencies, easy to install

## Prerequisites

- Go 1.26+ (for building from source)
- A [Statuscast](https://statuscast.com) account with API access

## Installation

### Download prebuilt binary

Download the latest release for your platform from [GitHub Releases](https://github.com/gunnaraasen/statuscast-mcp-server/releases).

```bash
# macOS arm64
curl -L https://github.com/gunnaraasen/statuscast-mcp-server/releases/latest/download/statuscast-mcp-server_Darwin_arm64.tar.gz | tar xz
chmod +x statuscast-mcp-server
mv statuscast-mcp-server /usr/local/bin/
```

### Install with go install

```bash
go install github.com/gunnaraasen/statuscast-mcp-server/cmd/server@latest
```

### Build from source

```bash
git clone https://github.com/gunnaraasen/statuscast-mcp-server.git
cd statuscast-mcp-server
make build
# Binary at: bin/statuscast-mcp-server
```

## Configuration

All configuration is via environment variables:

| Variable              | Required | Default | Description                                          |
|-----------------------|----------|---------|------------------------------------------------------|
| `STATUSCAST_TOKEN`    | Yes      | —       | Non-expiring account-level API token                 |
| `STATUSCAST_DOMAIN`   | Yes      | —       | Your Statuscast domain (e.g. `myco.statuscast.com`)  |
| `TRANSPORT`           | No       | `stdio` | Transport mode: `stdio` or `http`                    |
| `PORT`                | No       | `8080`  | HTTP port (only used when `TRANSPORT=http`)           |

## Claude Desktop Configuration

Add to your `claude_desktop_config.json` (`~/Library/Application Support/Claude/claude_desktop_config.json` on macOS):

```json
{
  "mcpServers": {
    "statuscast": {
      "command": "/usr/local/bin/statuscast-mcp-server",
      "env": {
        "STATUSCAST_TOKEN": "your-api-token",
        "STATUSCAST_DOMAIN": "myco.statuscast.com"
      }
    }
  }
}
```

## Kiro Configuration

Add to `.kiro/settings/mcp.json` in your project:

```json
{
  "mcpServers": {
    "statuscast": {
      "command": "statuscast-mcp-server",
      "env": {
        "STATUSCAST_TOKEN": "your-api-token",
        "STATUSCAST_DOMAIN": "myco.statuscast.com"
      }
    }
  }
}
```

## HTTP Transport

For remote deployments, run with `TRANSPORT=http`:

```bash
STATUSCAST_TOKEN=your-token STATUSCAST_DOMAIN=myco.statuscast.com TRANSPORT=http PORT=8080 statuscast-mcp-server
```

Configure your MCP client to connect to `http://your-host:8080/mcp`.

## Tool Reference

| Tool | Description | Key Parameters |
|------|-------------|----------------|
| `create_incident` | Create a new incident | `subject`*, `message`*, `type`, `affected_components`, `auto_publish`, `auto_close` |
| `get_incident` | Get an incident by ID | `incident_id`* |
| `update_incident` | Update an existing incident | `incident_id`*, `type`, `subject`, `message`, `status` |
| `search_incidents` | Search and list incidents | `query`, `status`, `limit` |
| `list_components` | List all components and current status | — |
| `get_component_history` | Get component uptime history | `component_id`, `time_range` |
| `create_subscriber` | Create a new subscriber | `email`*, `components`, `audience_groups` |
| `list_content_templates` | List all content templates | — |
| `create_content_template` | Create a content template | `name`*, `subject`*, `message`*, `type` |

*Required parameter

### `time_range` values for `get_component_history`

`Last7Days`, `Last30Days`, `ThisMonth`, `LastMonth`

## Development

```bash
make build      # Build binary to bin/statuscast-mcp-server
make test       # Run tests
make lint       # Run golangci-lint
make run-stdio  # Run with stdio transport
make run-http   # Run with HTTP transport on :8080
```

## Contributing

Contributions are welcome. Please open an issue first to discuss significant changes.

## License

MIT
