package config

import (
	"os"
	"testing"
)

func TestLoad_MissingToken(t *testing.T) {
	os.Unsetenv("STATUSCAST_TOKEN")
	os.Setenv("STATUSCAST_DOMAIN", "test.statuscast.com")
	t.Cleanup(func() { os.Unsetenv("STATUSCAST_DOMAIN") })

	_, err := Load()
	if err == nil {
		t.Fatal("expected error when STATUSCAST_TOKEN is missing")
	}
}

func TestLoad_MissingDomain(t *testing.T) {
	os.Setenv("STATUSCAST_TOKEN", "tok123")
	os.Unsetenv("STATUSCAST_DOMAIN")
	t.Cleanup(func() { os.Unsetenv("STATUSCAST_TOKEN") })

	_, err := Load()
	if err == nil {
		t.Fatal("expected error when STATUSCAST_DOMAIN is missing")
	}
}

func TestLoad_Defaults(t *testing.T) {
	os.Setenv("STATUSCAST_TOKEN", "tok123")
	os.Setenv("STATUSCAST_DOMAIN", "test.statuscast.com")
	os.Unsetenv("TRANSPORT")
	os.Unsetenv("PORT")
	t.Cleanup(func() {
		os.Unsetenv("STATUSCAST_TOKEN")
		os.Unsetenv("STATUSCAST_DOMAIN")
	})

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
	os.Setenv("STATUSCAST_TOKEN", "tok123")
	os.Setenv("STATUSCAST_DOMAIN", "test.statuscast.com")
	os.Setenv("TRANSPORT", "http")
	os.Setenv("PORT", "9090")
	t.Cleanup(func() {
		os.Unsetenv("STATUSCAST_TOKEN")
		os.Unsetenv("STATUSCAST_DOMAIN")
		os.Unsetenv("TRANSPORT")
		os.Unsetenv("PORT")
	})

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
