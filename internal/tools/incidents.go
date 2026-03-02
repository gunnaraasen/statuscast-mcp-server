package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gunnaraasen/statuscast-mcp-server/internal/client"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type incidentGroupArg struct {
	ID     int    `json:"id"               jsonschema:"Group ID (required)"`
	Action string `json:"action,omitempty" jsonschema:"Action: Add or Remove (defaults to Add)"`
}

type createIncidentArgs struct {
	Subject            string             `json:"subject"                        jsonschema:"The incident subject/title (required)"`
	Message            string             `json:"message"                        jsonschema:"The incident message body (required)"`
	IncidentType       int                `json:"incident_type,omitempty"        jsonschema:"Incident type integer (1=ServiceUnavailable, 2=ScheduledMaintenance, 3=Informational; mapping unverified)"`
	AffectedComponents []int              `json:"affected_components,omitempty"  jsonschema:"Component IDs affected by this incident"`
	SendNotifications  bool               `json:"send_notifications,omitempty"   jsonschema:"Notify subscribers immediately when true"`
	Active             bool               `json:"active,omitempty"               jsonschema:"Set to true to make the incident active/published"`
	DateToPost         string             `json:"date_to_post,omitempty"         jsonschema:"Scheduled post date/time (ISO 8601)"`
	HappeningNow       bool               `json:"happening_now,omitempty"        jsonschema:"Set to true if the incident is happening now"`
	TreatAsDowntime    bool               `json:"treat_as_downtime,omitempty"    jsonschema:"Set to true to count this incident toward downtime metrics"`
	EstimatedDuration  int                `json:"estimated_duration,omitempty"   jsonschema:"Estimated duration in minutes"`
	Groups             []incidentGroupArg `json:"groups,omitempty"               jsonschema:"Component groups to associate with this incident"`
	ProviderIncidentID string             `json:"provider_incident_id,omitempty" jsonschema:"External provider incident ID for cross-referencing"`
}

func createIncidentHandler(c *client.Client) mcp.ToolHandlerFor[createIncidentArgs, any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, args createIncidentArgs) (*mcp.CallToolResult, any, error) {
		groups := make([]client.IncidentGroup, len(args.Groups))
		for i, g := range args.Groups {
			groups[i] = client.IncidentGroup{ID: g.ID, Action: g.Action}
		}
		incident, err := c.CreateIncident(ctx, client.CreateIncidentRequest{
			Subject:            args.Subject,
			Message:            args.Message,
			IncidentType:       args.IncidentType,
			AffectedComponents: args.AffectedComponents,
			SendNotifications:  args.SendNotifications,
			Active:             args.Active,
			DateToPost:         args.DateToPost,
			HappeningNow:       args.HappeningNow,
			TreatAsDowntime:    args.TreatAsDowntime,
			EstimatedDuration:  args.EstimatedDuration,
			Groups:             groups,
			ProviderIncidentID: args.ProviderIncidentID,
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
	IncidentID         int                `json:"incident_id"                    jsonschema:"The ID of the incident to update (required)"`
	IncidentType       int                `json:"incident_type,omitempty"        jsonschema:"Updated incident type integer (1=ServiceUnavailable, 2=ScheduledMaintenance, 3=Informational; mapping unverified)"`
	Subject            string             `json:"subject,omitempty"              jsonschema:"Incident subject/title (auto-fetched from the API if omitted)"`
	Message            string             `json:"message,omitempty"              jsonschema:"Incident message body (auto-fetched from the API if omitted)"`
	Active             *bool              `json:"active,omitempty"               jsonschema:"Set to true to activate or false to deactivate/resolve the incident"`
	AffectedComponents []int              `json:"affected_components,omitempty"  jsonschema:"Updated list of affected component IDs"`
	DateToPost         string             `json:"date_to_post,omitempty"         jsonschema:"Scheduled post date/time (ISO 8601)"`
	HappeningNow       bool               `json:"happening_now,omitempty"        jsonschema:"Set to true if the incident is happening now"`
	TreatAsDowntime    bool               `json:"treat_as_downtime,omitempty"    jsonschema:"Set to true to count this incident toward downtime metrics"`
	EstimatedDuration  int                `json:"estimated_duration,omitempty"   jsonschema:"Estimated duration in minutes"`
	Groups             []incidentGroupArg `json:"groups,omitempty"               jsonschema:"Component groups to associate with this incident"`
	ProviderIncidentID string             `json:"provider_incident_id,omitempty" jsonschema:"External provider incident ID for cross-referencing"`
}

func updateIncidentHandler(c *client.Client) mcp.ToolHandlerFor[updateIncidentArgs, any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, args updateIncidentArgs) (*mcp.CallToolResult, any, error) {
		subject, message := args.Subject, args.Message
		if subject == "" || message == "" {
			current, err := c.GetIncident(ctx, args.IncidentID)
			if err != nil {
				return nil, nil, fmt.Errorf("auto-fetching incident %d for required subject/message: %w", args.IncidentID, err)
			}
			if subject == "" {
				subject = current.Subject
			}
			if message == "" {
				message = current.Message
				if message == "" {
					message = current.Body // GET may use "message" key
				}
			}
		}
		if subject == "" {
			return nil, nil, fmt.Errorf("subject is required by the API and could not be fetched; provide it explicitly")
		}

		groups := make([]client.IncidentGroup, len(args.Groups))
		for i, g := range args.Groups {
			groups[i] = client.IncidentGroup{ID: g.ID, Action: g.Action}
		}

		incident, err := c.UpdateIncident(ctx, client.UpdateIncidentRequest{
			ID:                 args.IncidentID,
			IncidentType:       args.IncidentType,
			Subject:            subject,
			Message:            message,
			Active:             args.Active,
			AffectedComponents: args.AffectedComponents,
			DateToPost:         args.DateToPost,
			HappeningNow:       args.HappeningNow,
			TreatAsDowntime:    args.TreatAsDowntime,
			EstimatedDuration:  args.EstimatedDuration,
			Groups:             groups,
			ProviderIncidentID: args.ProviderIncidentID,
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
	Sorting    string `json:"sorting,omitempty"     jsonschema:"Sort order: Ascending or Descending"`
}

func searchIncidentsHandler(c *client.Client) mcp.ToolHandlerFor[searchIncidentsArgs, any] {
	return func(ctx context.Context, _ *mcp.CallToolRequest, args searchIncidentsArgs) (*mcp.CallToolResult, any, error) {
		resp, err := c.SearchIncidents(ctx, client.SearchIncidentsRequest{
			TextSearch: args.TextSearch,
			PageNumber: args.PageNumber,
			PageSize:   args.PageSize,
			Sorting:    args.Sorting,
		})
		if err != nil {
			return nil, nil, err
		}
		data, _ := json.Marshal(resp)
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
