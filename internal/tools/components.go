package tools

import (
	"context"
	"encoding/json"

	"github.com/gunnaraasen/statuscast-mcp-server/internal/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type listComponentsArgs struct{}

func listComponentsHandler(c *client.Client) mcp.ToolHandlerFor[listComponentsArgs, any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, _ listComponentsArgs) (*mcp.CallToolResult, any, error) {
		components, err := c.ListComponents(ctx)
		if err != nil {
			return nil, nil, err
		}
		data, _ := json.Marshal(components)
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
		}, nil, nil
	}
}

type getComponentHistoryArgs struct {
	ComponentID int    `json:"component_id" jsonschema:"ID of the component to get status history for (required)"`
	Range_      string `json:"range,omitempty" jsonschema:"Time range for history: ThisWeek, ThisMonth, ThisYear, LastWeek, LastMonth, LastYear, Last7Days, Last30Days, Last60Days"`
}

func getComponentHistoryHandler(c *client.Client) mcp.ToolHandlerFor[getComponentHistoryArgs, any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, args getComponentHistoryArgs) (*mcp.CallToolResult, any, error) {
		history, err := c.GetComponentHistory(ctx, args.ComponentID, args.Range_)
		if err != nil {
			return nil, nil, err
		}
		data, _ := json.Marshal(history)
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
		}, nil, nil
	}
}

// RegisterComponentTools registers all component-related tools with the server.
func RegisterComponentTools(s *mcp.Server, c *client.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_components",
		Description: "List all Statuscast components and their current status",
	}, listComponentsHandler(c))

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_component_history",
		Description: "Get historical uptime and status data for a Statuscast component",
	}, getComponentHistoryHandler(c))
}
