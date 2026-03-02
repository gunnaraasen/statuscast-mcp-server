package tools

import (
	"github.com/gunnaraasen/statuscast-mcp-server/internal/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegisterAll registers all tool handlers with the server.
func RegisterAll(s *mcp.Server, c *client.Client) {
	RegisterIncidentTools(s, c)
	RegisterComponentTools(s, c)
	RegisterSubscriberTools(s, c)
	RegisterTemplateTools(s, c)
}
