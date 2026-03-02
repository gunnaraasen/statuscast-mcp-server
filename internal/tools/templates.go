package tools

import (
	"context"
	"encoding/json"

	"github.com/gunnaraasen/statuscast-mcp-server/internal/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type listContentTemplatesArgs struct{}

func listContentTemplatesHandler(c *client.Client) mcp.ToolHandlerFor[listContentTemplatesArgs, any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, _ listContentTemplatesArgs) (*mcp.CallToolResult, any, error) {
		templates, err := c.ListContentTemplates(ctx)
		if err != nil {
			return nil, nil, err
		}
		data, _ := json.Marshal(templates)
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
		}, nil, nil
	}
}

type createContentTemplateArgs struct {
	Name    string `json:"name"           jsonschema:"Template name for identification (required)"`
	Subject string `json:"subject"        jsonschema:"Template subject line (required)"`
	Message string `json:"message"        jsonschema:"Template message body (required)"`
	Type    string `json:"type,omitempty" jsonschema:"Template type (e.g. incident, maintenance)"`
}

func createContentTemplateHandler(c *client.Client) mcp.ToolHandlerFor[createContentTemplateArgs, any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, args createContentTemplateArgs) (*mcp.CallToolResult, any, error) {
		template, err := c.CreateContentTemplate(ctx, client.CreateContentTemplateRequest{
			Name:    args.Name,
			Subject: args.Subject,
			Message: args.Message,
			Type:    args.Type,
		})
		if err != nil {
			return nil, nil, err
		}
		data, _ := json.Marshal(template)
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
		}, nil, nil
	}
}

// RegisterTemplateTools registers all content template tools with the server.
func RegisterTemplateTools(s *mcp.Server, c *client.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "list_content_templates",
		Description: "List all Statuscast content templates",
	}, listContentTemplatesHandler(c))

	mcp.AddTool(s, &mcp.Tool{
		Name:        "create_content_template",
		Description: "Create a new Statuscast content template for reuse across incidents",
	}, createContentTemplateHandler(c))
}
