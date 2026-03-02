package tools

import (
	"context"
	"encoding/json"

	"github.com/gunnaraasen/statuscast-mcp-server/internal/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type createSubscriberArgs struct {
	Email          string `json:"email"                     jsonschema:"Email address of the subscriber (required)"`
	Components     []int  `json:"components,omitempty"      jsonschema:"Component IDs the subscriber wants to receive notifications for"`
	AudienceGroups []int  `json:"audience_groups,omitempty" jsonschema:"Audience group IDs to add the subscriber to"`
}

func createSubscriberHandler(c *client.Client) mcp.ToolHandlerFor[createSubscriberArgs, any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, args createSubscriberArgs) (*mcp.CallToolResult, any, error) {
		subscriber, err := c.CreateSubscriber(ctx, client.CreateSubscriberRequest{
			Email:          args.Email,
			Components:     args.Components,
			AudienceGroups: args.AudienceGroups,
		})
		if err != nil {
			return nil, nil, err
		}
		data, _ := json.Marshal(subscriber)
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
		}, nil, nil
	}
}

// RegisterSubscriberTools registers all subscriber-related tools with the server.
func RegisterSubscriberTools(s *mcp.Server, c *client.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "create_subscriber",
		Description: "Create a new Statuscast subscriber to receive status page notifications",
	}, createSubscriberHandler(c))
}
