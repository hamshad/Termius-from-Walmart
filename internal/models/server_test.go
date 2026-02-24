package models

import "testing"

func TestServerFilterValue(t *testing.T) {
	s := Server{Name: "My VPS"}
	if s.FilterValue() != "My VPS" {
		t.Errorf("FilterValue() = %q, want %q", s.FilterValue(), "My VPS")
	}
}

func TestServerTitle(t *testing.T) {
	s := Server{Name: "Production"}
	if s.Title() != "Production" {
		t.Errorf("Title() = %q, want %q", s.Title(), "Production")
	}
}

func TestServerDescription(t *testing.T) {
	s := Server{Username: "root", Host: "10.0.0.1", Port: 22}
	expected := "root@10.0.0.1:22"
	if s.Description() != expected {
		t.Errorf("Description() = %q, want %q", s.Description(), expected)
	}
}

func TestServerDescriptionCustomPort(t *testing.T) {
	s := Server{Username: "admin", Host: "example.com", Port: 2222}
	expected := "admin@example.com:2222"
	if s.Description() != expected {
		t.Errorf("Description() = %q, want %q", s.Description(), expected)
	}
}

func TestNewConfig(t *testing.T) {
	cfg := NewConfig()
	if cfg == nil {
		t.Fatal("NewConfig() returned nil")
	}
	if cfg.NextID != 1 {
		t.Errorf("NextID = %d, want 1", cfg.NextID)
	}
	if len(cfg.Servers) != 0 {
		t.Errorf("Servers length = %d, want 0", len(cfg.Servers))
	}
}
