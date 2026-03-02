package config

import (
	"testing"
)

func TestLoad_MissingToken(t *testing.T) {
	t.Setenv("STATUSCAST_DOMAIN", "test.statuscast.com")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error when STATUSCAST_TOKEN is missing")
	}
}

func TestLoad_MissingDomain(t *testing.T) {
	t.Setenv("STATUSCAST_TOKEN", "tok123")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error when STATUSCAST_DOMAIN is missing")
	}
}

func TestLoad_Defaults(t *testing.T) {
	t.Setenv("STATUSCAST_TOKEN", "tok123")
	t.Setenv("STATUSCAST_DOMAIN", "test.statuscast.com")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Transport != "stdio" {
		t.Errorf("expected Transport=stdio, got %q", cfg.Transport)
	}
	if cfg.Port != "8080" {
		t.Errorf("expected Port=8080, got %q", cfg.Port)
	}
}

func TestLoad_CustomTransport(t *testing.T) {
	t.Setenv("STATUSCAST_TOKEN", "tok123")
	t.Setenv("STATUSCAST_DOMAIN", "test.statuscast.com")
	t.Setenv("TRANSPORT", "http")
	t.Setenv("PORT", "9090")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Transport != "http" {
		t.Errorf("expected Transport=http, got %q", cfg.Transport)
	}
	if cfg.Port != "9090" {
		t.Errorf("expected Port=9090, got %q", cfg.Port)
	}
	if cfg.Token != "tok123" {
		t.Errorf("expected Token=tok123, got %q", cfg.Token)
	}
	if cfg.Domain != "test.statuscast.com" {
		t.Errorf("expected Domain=test.statuscast.com, got %q", cfg.Domain)
	}
}
