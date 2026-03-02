package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
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
	ID                 string   `json:"id"`
	Subject            string   `json:"subject"`
	Message            string   `json:"message"`
	Type               string   `json:"type"`
	Status             string   `json:"status"`
	AffectedComponents []string `json:"affected_components,omitempty"`
	AutoPublish        bool     `json:"auto_publish,omitempty"`
	AutoClose          bool     `json:"auto_close,omitempty"`
	CreatedAt          string   `json:"created_at,omitempty"`
	UpdatedAt          string   `json:"updated_at,omitempty"`
}

// Component represents a Statuscast component.
type Component struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Status      string `json:"status"`
	Description string `json:"description,omitempty"`
	GroupID     string `json:"group_id,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}

// ComponentHistory represents the historical status of a component.
type ComponentHistory struct {
	ComponentID string            `json:"component_id,omitempty"`
	TimeRange   string            `json:"time_range,omitempty"`
	Uptime      float64           `json:"uptime,omitempty"`
	History     []HistoryEntry    `json:"history,omitempty"`
}

// HistoryEntry represents a single point in component history.
type HistoryEntry struct {
	Date   string `json:"date"`
	Status string `json:"status"`
}

// Subscriber represents a Statuscast subscriber.
type Subscriber struct {
	ID              string   `json:"id"`
	Email           string   `json:"email"`
	Components      []string `json:"components,omitempty"`
	AudienceGroups  []string `json:"audience_groups,omitempty"`
	CreatedAt       string   `json:"created_at,omitempty"`
}

// ContentTemplate represents a Statuscast content template.
type ContentTemplate struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Subject   string `json:"subject"`
	Message   string `json:"message"`
	Type      string `json:"type,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

// --- Request types ---

// CreateIncidentRequest holds the parameters for creating an incident.
type CreateIncidentRequest struct {
	Subject            string   `json:"subject"`
	Message            string   `json:"message"`
	Type               string   `json:"type,omitempty"`
	AffectedComponents []string `json:"affected_components,omitempty"`
	AutoPublish        bool     `json:"auto_publish,omitempty"`
	AutoClose          bool     `json:"auto_close,omitempty"`
}

// UpdateIncidentRequest holds the parameters for updating an incident.
type UpdateIncidentRequest struct {
	ID      string `json:"-"`
	Type    string `json:"type,omitempty"`
	Subject string `json:"subject,omitempty"`
	Message string `json:"message,omitempty"`
	Status  string `json:"status,omitempty"`
}

// SearchIncidentsRequest holds the parameters for searching incidents.
type SearchIncidentsRequest struct {
	Query  string
	Status string
	Limit  int
}

// CreateSubscriberRequest holds the parameters for creating a subscriber.
type CreateSubscriberRequest struct {
	Email          string   `json:"email"`
	Components     []string `json:"components,omitempty"`
	AudienceGroups []string `json:"audience_groups,omitempty"`
}

// CreateContentTemplateRequest holds the parameters for creating a content template.
type CreateContentTemplateRequest struct {
	Name    string `json:"name"`
	Subject string `json:"subject"`
	Message string `json:"message"`
	Type    string `json:"type,omitempty"`
}

// --- API methods ---

// CreateIncident creates a new incident.
func (c *Client) CreateIncident(ctx context.Context, req CreateIncidentRequest) (*Incident, error) {
	var incident Incident
	if err := c.do(ctx, http.MethodPost, "/incidents", req, &incident); err != nil {
		return nil, err
	}
	return &incident, nil
}

// GetIncident retrieves an incident by ID.
func (c *Client) GetIncident(ctx context.Context, id string) (*Incident, error) {
	var incident Incident
	if err := c.do(ctx, http.MethodGet, "/incidents/"+id, nil, &incident); err != nil {
		return nil, err
	}
	return &incident, nil
}

// UpdateIncident updates an existing incident.
func (c *Client) UpdateIncident(ctx context.Context, req UpdateIncidentRequest) (*Incident, error) {
	var incident Incident
	if err := c.do(ctx, http.MethodPatch, "/incidents/"+req.ID, req, &incident); err != nil {
		return nil, err
	}
	return &incident, nil
}

// SearchIncidents searches and lists incidents with optional filters.
func (c *Client) SearchIncidents(ctx context.Context, req SearchIncidentsRequest) ([]Incident, error) {
	params := url.Values{}
	if req.Query != "" {
		params.Set("search", req.Query)
	}
	if req.Status != "" {
		params.Set("status", req.Status)
	}
	if req.Limit > 0 {
		params.Set("limit", strconv.Itoa(req.Limit))
	}

	path := "/incidents"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	var incidents []Incident
	if err := c.do(ctx, http.MethodGet, path, nil, &incidents); err != nil {
		return nil, err
	}
	return incidents, nil
}

// ListComponents returns all components.
func (c *Client) ListComponents(ctx context.Context) ([]Component, error) {
	var components []Component
	if err := c.do(ctx, http.MethodGet, "/components", nil, &components); err != nil {
		return nil, err
	}
	return components, nil
}

// GetComponentHistory returns historical status data for a component.
func (c *Client) GetComponentHistory(ctx context.Context, id, timeRange string) (*ComponentHistory, error) {
	params := url.Values{}
	if timeRange != "" {
		params.Set("time_range", timeRange)
	}

	path := "/components"
	if id != "" {
		path += "/" + id
	}
	path += "/history"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	var history ComponentHistory
	if err := c.do(ctx, http.MethodGet, path, nil, &history); err != nil {
		return nil, err
	}
	return &history, nil
}

// CreateSubscriber creates a new subscriber.
func (c *Client) CreateSubscriber(ctx context.Context, req CreateSubscriberRequest) (*Subscriber, error) {
	var subscriber Subscriber
	if err := c.do(ctx, http.MethodPost, "/subscribers", req, &subscriber); err != nil {
		return nil, err
	}
	return &subscriber, nil
}

// ListContentTemplates returns all content templates.
func (c *Client) ListContentTemplates(ctx context.Context) ([]ContentTemplate, error) {
	var templates []ContentTemplate
	if err := c.do(ctx, http.MethodGet, "/content-templates", nil, &templates); err != nil {
		return nil, err
	}
	return templates, nil
}

// CreateContentTemplate creates a new content template.
func (c *Client) CreateContentTemplate(ctx context.Context, req CreateContentTemplateRequest) (*ContentTemplate, error) {
	var template ContentTemplate
	if err := c.do(ctx, http.MethodPost, "/content-templates", req, &template); err != nil {
		return nil, err
	}
	return &template, nil
}
