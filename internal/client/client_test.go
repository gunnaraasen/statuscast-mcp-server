package client

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestClient(server *httptest.Server) *Client {
	c := New("placeholder.statuscast.com", "test-token")
	c.baseURL = server.URL
	return c
}

func TestDo_AuthHeader(t *testing.T) {
	var gotAuth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}"))
	}))
	defer srv.Close()

	c := newTestClient(srv)
	c.do(context.Background(), http.MethodGet, "/test", nil, nil) //nolint:errcheck

	if gotAuth != "Bearer test-token" {
		t.Errorf("expected 'Bearer test-token', got %q", gotAuth)
	}
}

func TestDo_APIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message":"invalid token"}`))
	}))
	defer srv.Close()

	c := newTestClient(srv)
	err := c.do(context.Background(), http.MethodGet, "/test", nil, nil)

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T: %v", err, err)
	}
	if apiErr.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", apiErr.StatusCode)
	}
	if apiErr.Message != "invalid token" {
		t.Errorf("expected message 'invalid token', got %q", apiErr.Message)
	}
}

func TestDo_UnmarshalResponse(t *testing.T) {
	incident := Incident{ID: "inc_123", Subject: "Test Incident"}
	body, _ := json.Marshal(incident)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}))
	defer srv.Close()

	c := newTestClient(srv)
	var got Incident
	if err := c.do(context.Background(), http.MethodGet, "/incidents/inc_123", nil, &got); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != "inc_123" {
		t.Errorf("expected ID=inc_123, got %q", got.ID)
	}
	if got.Subject != "Test Incident" {
		t.Errorf("expected Subject='Test Incident', got %q", got.Subject)
	}
}

func TestDo_JSONRequestBody(t *testing.T) {
	var gotBody string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("{}"))
	}))
	defer srv.Close()

	c := newTestClient(srv)
	req := CreateIncidentRequest{Subject: "Outage", Message: "Services down"}
	c.do(context.Background(), http.MethodPost, "/incidents", req, nil) //nolint:errcheck

	var parsed map[string]any
	if err := json.Unmarshal([]byte(gotBody), &parsed); err != nil {
		t.Fatalf("request body is not valid JSON: %v", err)
	}
	if parsed["subject"] != "Outage" {
		t.Errorf("expected subject='Outage', got %v", parsed["subject"])
	}
}

func TestAPIError_Error(t *testing.T) {
	err := &APIError{StatusCode: 404, Message: "not found"}
	want := "statuscast API error 404: not found"
	if err.Error() != want {
		t.Errorf("expected %q, got %q", want, err.Error())
	}
}
