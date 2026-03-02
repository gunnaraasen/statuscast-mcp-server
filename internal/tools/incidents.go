package tools

import (
	"context"
	"encoding/json"

	"github.com/gunnaraasen/statuscast-mcp-server/internal/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type createIncidentArgs struct {
	Subject            string `json:"subject"                        jsonschema:"The incident subject/title (required)"`
	Message            string `json:"message"                        jsonschema:"The incident message body (required)"`
	IncidentType       int    `json:"incident_type,omitempty"        jsonschema:"Incident type: 1=ServiceUnavailable, 2=ScheduledMaintenance, 3=Informational"`
	AffectedComponents []int  `json:"affected_components,omitempty"  jsonschema:"Component IDs affected by this incident"`
	SendNotifications  bool   `json:"send_notifications,omitempty"   jsonschema:"Notify subscribers immediately when true"`
	Active             bool   `json:"active,omitempty"               jsonschema:"Set to true to make the incident active/published"`
}

func createIncidentHandler(c *client.Client) mcp.ToolHandlerFor[createIncidentArgs, any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, args createIncidentArgs) (*mcp.CallToolResult, any, error) {
		incident, err := c.CreateIncident(ctx, client.CreateIncidentRequest{
			Subject:            args.Subject,
			Message:            args.Message,
			IncidentType:       args.IncidentType,
			AffectedComponents: args.AffectedComponents,
			SendNotifications:  args.SendNotifications,
			Active:             args.Active,
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
	IncidentID int `json:"incident_id" jsonschema:"The ID of the incident to retrieve (required)"`
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
	IncidentID         int    `json:"incident_id"                    jsonschema:"The ID of the incident to update (required)"`
	IncidentType       int    `json:"incident_type,omitempty"        jsonschema:"Updated incident type: 1=ServiceUnavailable, 2=ScheduledMaintenance, 3=Informational"`
	Subject            string `json:"subject"                        jsonschema:"Incident subject/title (required by the API)"`
	Message            string `json:"message"                        jsonschema:"Incident message body (required by the API)"`
	Active             *bool  `json:"active,omitempty"               jsonschema:"Set to true to activate or false to deactivate/resolve the incident"`
	AffectedComponents []int  `json:"affected_components,omitempty"  jsonschema:"Updated list of affected component IDs"`
}

func updateIncidentHandler(c *client.Client) mcp.ToolHandlerFor[updateIncidentArgs, any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, args updateIncidentArgs) (*mcp.CallToolResult, any, error) {
		incident, err := c.UpdateIncident(ctx, client.UpdateIncidentRequest{
			ID:                 args.IncidentID,
			IncidentType:       args.IncidentType,
			Subject:            args.Subject,
			Message:            args.Message,
			Active:             args.Active,
			AffectedComponents: args.AffectedComponents,
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
	TextSearch string `json:"text_search,omitempty" jsonschema:"Search query string to filter incidents by subject or message"`
	PageNumber int    `json:"page_number,omitempty" jsonschema:"Page number for pagination (1-based)"`
	PageSize   int    `json:"page_size,omitempty"   jsonschema:"Number of incidents per page"`
}

func searchIncidentsHandler(c *client.Client) mcp.ToolHandlerFor[searchIncidentsArgs, any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, args searchIncidentsArgs) (*mcp.CallToolResult, any, error) {
		incidents, err := c.SearchIncidents(ctx, client.SearchIncidentsRequest{
			TextSearch: args.TextSearch,
			PageNumber: args.PageNumber,
			PageSize:   args.PageSize,
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
		Description: "Search and list Statuscast incidents with optional text search and pagination",
	}, searchIncidentsHandler(c))
}
