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
	Event      string `json:"event"                  jsonschema:"Lifecycle event: NewIncident, UpdateIncident, ResolveIncident, RCAIncident, InvestigatingUpdate, MonitoringUpdate, IdentifiedUpdate (required)"`
	Status     string `json:"status,omitempty"       jsonschema:"Component status: Available, Unavailable, DegradedPerformance, Maintenance, Investigating, Monitoring, Identified"`
	PostType   string `json:"post_type,omitempty"    jsonschema:"Update type: Informational, Closed, RootCause, Investigating, Monitoring, Identified, Migration"`
	Subject    string `json:"subject,omitempty"      jsonschema:"Template subject line"`
	Contents   string `json:"contents,omitempty"     jsonschema:"Template message body"`
	Components []int  `json:"components,omitempty"   jsonschema:"Component IDs this template applies to"`
	Groups     []int  `json:"groups,omitempty"       jsonschema:"Group IDs this template applies to"`
}

func createContentTemplateHandler(c *client.Client) mcp.ToolHandlerFor[createContentTemplateArgs, any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, args createContentTemplateArgs) (*mcp.CallToolResult, any, error) {
		template, err := c.CreateContentTemplate(ctx, client.CreateContentTemplateRequest{
			Event:      args.Event,
			Status:     args.Status,
			PostType:   args.PostType,
			Subject:    args.Subject,
			Contents:   args.Contents,
			Components: args.Components,
			Groups:     args.Groups,
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
