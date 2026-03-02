package tools

import (
	"context"
	"encoding/json"

	"github.com/gunnaraasen/statuscast-mcp-server/internal/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type createIncidentArgs struct {
	Subject            string   `json:"subject"             jsonschema:"The incident subject/title (required)"`
	Message            string   `json:"message"             jsonschema:"The incident message body (required)"`
	Type               string   `json:"type,omitempty"      jsonschema:"Incident type (e.g. investigating, identified, monitoring, resolved)"`
	AffectedComponents []string `json:"affected_components,omitempty" jsonschema:"Component IDs affected by this incident"`
	AutoPublish        bool     `json:"auto_publish,omitempty" jsonschema:"Notify subscribers immediately when true"`
	AutoClose          bool     `json:"auto_close,omitempty"   jsonschema:"Auto-close incident when monitoring resolves when true"`
}

func createIncidentHandler(c *client.Client) mcp.ToolHandlerFor[createIncidentArgs, any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, args createIncidentArgs) (*mcp.CallToolResult, any, error) {
		incident, err := c.CreateIncident(ctx, client.CreateIncidentRequest{
			Subject:            args.Subject,
			Message:            args.Message,
			Type:               args.Type,
			AffectedComponents: args.AffectedComponents,
			AutoPublish:        args.AutoPublish,
			AutoClose:          args.AutoClose,
		})
		if err != nil {
			return nil, nil, err
		}
		data, _ := json.Marshal(incident)
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
		}, nil, nil
	}
}

type getIncidentArgs struct {
	IncidentID string `json:"incident_id" jsonschema:"The ID of the incident to retrieve (required)"`
}

func getIncidentHandler(c *client.Client) mcp.ToolHandlerFor[getIncidentArgs, any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, args getIncidentArgs) (*mcp.CallToolResult, any, error) {
		incident, err := c.GetIncident(ctx, args.IncidentID)
		if err != nil {
			return nil, nil, err
		}
		data, _ := json.Marshal(incident)
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
		}, nil, nil
	}
}

type updateIncidentArgs struct {
	IncidentID string `json:"incident_id"      jsonschema:"The ID of the incident to update (required)"`
	Type       string `json:"type,omitempty"   jsonschema:"Updated incident type"`
	Subject    string `json:"subject,omitempty" jsonschema:"Updated incident subject/title"`
	Message    string `json:"message,omitempty" jsonschema:"Updated incident message body"`
	Status     string `json:"status,omitempty" jsonschema:"Updated incident status (e.g. investigating, resolved)"`
}

func updateIncidentHandler(c *client.Client) mcp.ToolHandlerFor[updateIncidentArgs, any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, args updateIncidentArgs) (*mcp.CallToolResult, any, error) {
		incident, err := c.UpdateIncident(ctx, client.UpdateIncidentRequest{
			ID:      args.IncidentID,
			Type:    args.Type,
			Subject: args.Subject,
			Message: args.Message,
			Status:  args.Status,
		})
		if err != nil {
			return nil, nil, err
		}
		data, _ := json.Marshal(incident)
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
		}, nil, nil
	}
}

type searchIncidentsArgs struct {
	Query  string `json:"query,omitempty"  jsonschema:"Search query string to filter incidents by subject or message"`
	Status string `json:"status,omitempty" jsonschema:"Filter by status (e.g. investigating, resolved)"`
	Limit  int    `json:"limit,omitempty"  jsonschema:"Maximum number of incidents to return"`
}

func searchIncidentsHandler(c *client.Client) mcp.ToolHandlerFor[searchIncidentsArgs, any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, args searchIncidentsArgs) (*mcp.CallToolResult, any, error) {
		incidents, err := c.SearchIncidents(ctx, client.SearchIncidentsRequest{
			Query:  args.Query,
			Status: args.Status,
			Limit:  args.Limit,
		})
		if err != nil {
			return nil, nil, err
		}
		data, _ := json.Marshal(incidents)
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
		}, nil, nil
	}
}

// RegisterIncidentTools registers all incident-related tools with the server.
func RegisterIncidentTools(s *mcp.Server, c *client.Client) {
	mcp.AddTool(s, &mcp.Tool{
		Name:        "create_incident",
		Description: "Create a new Statuscast incident to notify subscribers of an outage or issue",
	}, createIncidentHandler(c))

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_incident",
		Description: "Get a Statuscast incident by ID",
	}, getIncidentHandler(c))

	mcp.AddTool(s, &mcp.Tool{
		Name:        "update_incident",
		Description: "Update an existing Statuscast incident (e.g. change status, add message update)",
	}, updateIncidentHandler(c))

	mcp.AddTool(s, &mcp.Tool{
		Name:        "search_incidents",
		Description: "Search and list Statuscast incidents with optional filters by query, status, or limit",
	}, searchIncidentsHandler(c))
}
