package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Client is an HTTP client for the Statuscast API v4.
type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

// New creates a new Statuscast API client.
func New(domain, token string) *Client {
	return &Client{
		baseURL:    "https://" + domain + "/api/v4",
		token:      token,
		httpClient: &http.Client{},
	}
}

// APIError represents a non-2xx response from the Statuscast API.
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("statuscast API error %d: %s", e.StatusCode, e.Message)
}

// do executes an HTTP request, marshaling body as JSON and unmarshaling the
// response into out. Pass nil for body when there is no request body.
func (c *Client) do(ctx context.Context, method, path string, body, out any) error {
	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshaling request body: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bodyReader)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Try to extract an error message from the response.
		var errResp struct {
			Message string `json:"message"`
			Error   string `json:"error"`
		}
		msg := http.StatusText(resp.StatusCode)
		if json.Unmarshal(respBody, &errResp) == nil {
			if errResp.Message != "" {
				msg = errResp.Message
			} else if errResp.Error != "" {
				msg = errResp.Error
			}
		}
		return &APIError{StatusCode: resp.StatusCode, Message: msg}
	}

	if out != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, out); err != nil {
			return fmt.Errorf("unmarshaling response: %w", err)
		}
	}

	return nil
}

// --- Data types ---

// Incident represents a Statuscast incident.
type Incident struct {
	ID                 int    `json:"id"`
	Subject            string `json:"messageSubject"`
	Message            string `json:"messageText"`
	IncidentType       int    `json:"incidentType,omitempty"`
	Active             bool   `json:"active,omitempty"`
	AffectedComponents []int  `json:"affectedComponents,omitempty"`
	SendNotifications  bool   `json:"sendNotifications,omitempty"`
}

// Component represents a Statuscast component.
type Component struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Status      string `json:"status"`
	Description string `json:"description,omitempty"`
	ParentID    int    `json:"parentId,omitempty"`
	IsHidden    bool   `json:"isHidden,omitempty"`
	ExternalID  string `json:"externalId,omitempty"`
}

// ComponentHistoryEntry represents a single status change event in component history.
type ComponentHistoryEntry struct {
	Status              string `json:"status"`
	DirectlyAffected    bool   `json:"directlyAffected"`
	CountTowardDowntime bool   `json:"countTowardDowntime"`
	DateChanged         string `json:"dateChanged"`
	IncidentID          *int   `json:"incidentId,omitempty"`
	ComponentID         int    `json:"componentId"`
}

// Subscriber represents a Statuscast subscriber.
type Subscriber struct {
	ID             int    `json:"id"`
	Email          string `json:"email"`
	Components     []int  `json:"components,omitempty"`
	AudienceGroups []int  `json:"audienceGroups,omitempty"`
	CreatedAt      string `json:"createdAt,omitempty"`
}

// ContentTemplate represents a Statuscast content template.
type ContentTemplate struct {
	ID         int    `json:"id"`
	Event      string `json:"event"`
	Status     string `json:"status,omitempty"`
	PostType   string `json:"postType,omitempty"`
	Components []int  `json:"components,omitempty"`
	Groups     []int  `json:"groups,omitempty"`
	Subject    string `json:"subject,omitempty"`
	Contents   string `json:"contents,omitempty"`
}

// --- Request types ---

// CreateIncidentRequest holds the parameters for creating an incident.
type CreateIncidentRequest struct {
	Subject            string `json:"messageSubject"`
	Message            string `json:"messageText"`
	IncidentType       int    `json:"incidentType,omitempty"`
	AffectedComponents []int  `json:"affectedComponents,omitempty"`
	SendNotifications  bool   `json:"sendNotifications,omitempty"`
	Active             bool   `json:"active,omitempty"`
}

// UpdateIncidentRequest holds the parameters for updating an incident.
// ID must be set; Active uses a pointer to distinguish false from omitted.
type UpdateIncidentRequest struct {
	ID                 int    `json:"id"`
	Subject            string `json:"messageSubject,omitempty"`
	Message            string `json:"messageText,omitempty"`
	IncidentType       int    `json:"incidentType,omitempty"`
	Active             *bool  `json:"active,omitempty"`
	AffectedComponents []int  `json:"affectedComponents,omitempty"`
}

// SearchIncidentsRequest holds the parameters for searching incidents via POST body.
type SearchIncidentsRequest struct {
	TextSearch string `json:"textSearch,omitempty"`
	PageNumber int    `json:"pageNumber,omitempty"`
	PageSize   int    `json:"pageSize,omitempty"`
}

// SearchIncidentsResponse wraps the paginated incident search result.
type SearchIncidentsResponse struct {
	Items      []Incident `json:"items"`
	TotalItems int        `json:"totalItems"`
	Pages      int        `json:"pages"`
}

// CreateSubscriberRequest holds the parameters for creating a subscriber.
type CreateSubscriberRequest struct {
	Email          string `json:"email"`
	Components     []int  `json:"components,omitempty"`
	AudienceGroups []int  `json:"audienceGroups,omitempty"`
}

// CreateContentTemplateRequest holds the parameters for creating a content template.
type CreateContentTemplateRequest struct {
	Event      string `json:"event"`
	Status     string `json:"status,omitempty"`
	PostType   string `json:"postType,omitempty"`
	Components []int  `json:"components,omitempty"`
	Groups     []int  `json:"groups,omitempty"`
	Subject    string `json:"subject,omitempty"`
	Contents   string `json:"contents,omitempty"`
}

// --- API methods ---

// CreateIncident creates a new incident.
func (c *Client) CreateIncident(ctx context.Context, req CreateIncidentRequest) (*Incident, error) {
	var incident Incident
	if err := c.do(ctx, http.MethodPost, "/incident", req, &incident); err != nil {
		return nil, err
	}
	return &incident, nil
}

// GetIncident retrieves an incident by ID.
func (c *Client) GetIncident(ctx context.Context, id int) (*Incident, error) {
	var incident Incident
	if err := c.do(ctx, http.MethodGet, fmt.Sprintf("/incident/%d", id), nil, &incident); err != nil {
		return nil, err
	}
	return &incident, nil
}

// UpdateIncident updates an existing incident. The ID must be set in the request body.
func (c *Client) UpdateIncident(ctx context.Context, req UpdateIncidentRequest) (*Incident, error) {
	var incident Incident
	if err := c.do(ctx, http.MethodPut, "/incident", req, &incident); err != nil {
		return nil, err
	}
	return &incident, nil
}

// SearchIncidents searches incidents via POST body and returns the result items.
func (c *Client) SearchIncidents(ctx context.Context, req SearchIncidentsRequest) ([]Incident, error) {
	var resp SearchIncidentsResponse
	if err := c.do(ctx, http.MethodPost, "/incidents", req, &resp); err != nil {
		return nil, err
	}
	return resp.Items, nil
}

// ListComponents returns all components.
func (c *Client) ListComponents(ctx context.Context) ([]Component, error) {
	var components []Component
	if err := c.do(ctx, http.MethodGet, "/components", nil, &components); err != nil {
		return nil, err
	}
	return components, nil
}

// GetComponentHistory returns historical status entries for a component.
// Pass id=0 to get history for all components. timeRange accepts values such as
// ThisWeek, ThisMonth, ThisYear, LastWeek, LastMonth, LastYear, Last7Days,
// Last30Days, Last60Days.
func (c *Client) GetComponentHistory(ctx context.Context, id int, timeRange string) ([]ComponentHistoryEntry, error) {
	params := url.Values{}
	if timeRange != "" {
		params.Set("range", timeRange)
	}

	var path string
	if id != 0 {
		path = fmt.Sprintf("/component/%d/history", id)
	} else {
		path = "/components/history"
	}
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	var history []ComponentHistoryEntry
	if err := c.do(ctx, http.MethodGet, path, nil, &history); err != nil {
		return nil, err
	}
	return history, nil
}

// CreateSubscriber creates a new subscriber.
func (c *Client) CreateSubscriber(ctx context.Context, req CreateSubscriberRequest) (*Subscriber, error) {
	var subscriber Subscriber
	if err := c.do(ctx, http.MethodPost, "/subscriber", req, &subscriber); err != nil {
		return nil, err
	}
	return &subscriber, nil
}

// ListContentTemplates returns all content templates.
func (c *Client) ListContentTemplates(ctx context.Context) ([]ContentTemplate, error) {
	var templates []ContentTemplate
	if err := c.do(ctx, http.MethodGet, "/contenttemplate", nil, &templates); err != nil {
		return nil, err
	}
	return templates, nil
}

// CreateContentTemplate creates a new content template.
func (c *Client) CreateContentTemplate(ctx context.Context, req CreateContentTemplateRequest) (*ContentTemplate, error) {
	var template ContentTemplate
	if err := c.do(ctx, http.MethodPost, "/contenttemplate", req, &template); err != nil {
		return nil, err
	}
	return &template, nil
}
