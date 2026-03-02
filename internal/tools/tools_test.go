package tools_test

import (
	"context"
	"testing"

	"github.com/gunnaraasen/statuscast-mcp-server/internal/client"
	"github.com/gunnaraasen/statuscast-mcp-server/internal/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// TestRegisterAll verifies that RegisterAll does not panic and registers the
// expected number of tools on the server.
func TestRegisterAll(t *testing.T) {
	s := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "0.0.1"}, nil)
	c := client.New("test.statuscast.com", "token")

	// Should not panic.
	tools.RegisterAll(s, c)

	// Verify all 9 tools are registered by listing them via in-memory transport.
	t1, t2 := mcp.NewInMemoryTransports()
	ctx := context.Background()

	_, err := s.Connect(ctx, t1, nil)
	if err != nil {
		t.Fatalf("Connect: %v", err)
	}

	mc := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "0.0.1"}, nil)
	session, err := mc.Connect(ctx, t2, nil)
	if err != nil {
		t.Fatalf("client Connect: %v", err)
	}
	defer session.Close()

	var count int
	for _, err := range session.Tools(ctx, nil) {
		if err != nil {
			t.Fatalf("Tools iteration: %v", err)
		}
		count++
	}

	const wantTools = 9
	if count != wantTools {
		t.Errorf("expected %d tools, got %d", wantTools, count)
	}
}

// TestToolNames verifies the registered tool names match expectations.
func TestToolNames(t *testing.T) {
	s := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "0.0.1"}, nil)
	c := client.New("test.statuscast.com", "token")
	tools.RegisterAll(s, c)

	t1, t2 := mcp.NewInMemoryTransports()
	ctx := context.Background()

	_, err := s.Connect(ctx, t1, nil)
	if err != nil {
		t.Fatalf("Connect: %v", err)
	}

	mc := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "0.0.1"}, nil)
	session, err := mc.Connect(ctx, t2, nil)
	if err != nil {
		t.Fatalf("client Connect: %v", err)
	}
	defer session.Close()

	want := map[string]bool{
		"create_incident":         true,
		"get_incident":            true,
		"update_incident":         true,
		"search_incidents":        true,
		"list_components":         true,
		"get_component_history":   true,
		"create_subscriber":       true,
		"list_content_templates":  true,
		"create_content_template": true,
	}

	for tool, err := range session.Tools(ctx, nil) {
		if err != nil {
			t.Fatalf("Tools iteration: %v", err)
		}
		if !want[tool.Name] {
			t.Errorf("unexpected tool: %q", tool.Name)
		}
		delete(want, tool.Name)
	}

	for name := range want {
		t.Errorf("missing expected tool: %q", name)
	}
}
