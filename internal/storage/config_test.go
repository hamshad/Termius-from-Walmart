package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"termius-from-walmart/internal/models"
)

func TestLoadConfig_NonexistentFile(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "nonexistent", "config.json")

	cfg := LoadConfig(path)
	if cfg == nil {
		t.Fatal("LoadConfig returned nil for nonexistent file")
	}
	if cfg.NextID != 1 {
		t.Errorf("NextID = %d, want 1", cfg.NextID)
	}
	if len(cfg.Servers) != 0 {
		t.Errorf("Servers length = %d, want 0", len(cfg.Servers))
	}
}

func TestLoadConfig_ValidFile(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "config.json")

	cfg := &models.Config{
		NextID: 5,
		Servers: []models.Server{
			{ID: 1, Name: "Test Server", Host: "10.0.0.1", Port: 22, Username: "root"},
			{ID: 2, Name: "Dev Server", Host: "dev.example.com", Port: 2222, Username: "dev"},
		},
	}

	data, _ := json.MarshalIndent(cfg, "", "  ")
	os.WriteFile(path, data, 0600)

	loaded := LoadConfig(path)
	if loaded.NextID != 5 {
		t.Errorf("NextID = %d, want 5", loaded.NextID)
	}
	if len(loaded.Servers) != 2 {
		t.Fatalf("Servers length = %d, want 2", len(loaded.Servers))
	}
	if loaded.Servers[0].Name != "Test Server" {
		t.Errorf("Servers[0].Name = %q, want %q", loaded.Servers[0].Name, "Test Server")
	}
	if loaded.Servers[1].Host != "dev.example.com" {
		t.Errorf("Servers[1].Host = %q, want %q", loaded.Servers[1].Host, "dev.example.com")
	}
}

func TestLoadConfig_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "config.json")

	os.WriteFile(path, []byte("not valid json {{{"), 0600)

	cfg := LoadConfig(path)
	if cfg == nil {
		t.Fatal("LoadConfig returned nil for invalid JSON")
	}
	if cfg.NextID != 1 {
		t.Errorf("NextID = %d, want 1 (should return default)", cfg.NextID)
	}
}

func TestSaveConfig(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "subdir", "config.json")

	cfg := &models.Config{
		NextID: 3,
		Servers: []models.Server{
			{ID: 1, Name: "Saved Server", Host: "192.168.1.1", Port: 22, Username: "admin"},
		},
	}

	err := SaveConfig(path, cfg)
	if err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	// Verify file exists
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Config file not found after save: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("File permissions = %o, want 0600", info.Mode().Perm())
	}

	// Verify content round-trips
	loaded := LoadConfig(path)
	if loaded.NextID != 3 {
		t.Errorf("NextID = %d, want 3", loaded.NextID)
	}
	if len(loaded.Servers) != 1 {
		t.Fatalf("Servers length = %d, want 1", len(loaded.Servers))
	}
	if loaded.Servers[0].Name != "Saved Server" {
		t.Errorf("Servers[0].Name = %q, want %q", loaded.Servers[0].Name, "Saved Server")
	}
}

func TestSaveConfig_CreatesDirectories(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "deep", "nested", "dir", "config.json")

	cfg := models.NewConfig()
	err := SaveConfig(path, cfg)
	if err != nil {
		t.Fatalf("SaveConfig failed to create nested dirs: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Errorf("Config file not found at deeply nested path: %v", err)
	}
}

func TestExportServers(t *testing.T) {
	tmpDir := t.TempDir()
	exportPath := filepath.Join(tmpDir, "export.json")

	servers := []models.Server{
		{ID: 1, Name: "Server A", Host: "10.0.0.1", Port: 22, Username: "root", Password: "pass123", SFTPPort: 22},
		{ID: 2, Name: "Server B", Host: "example.com", Port: 2222, Username: "user", PemKey: "-----BEGIN RSA PRIVATE KEY-----\ntest\n-----END RSA PRIVATE KEY-----", SFTPPort: 2222},
	}

	err := ExportServers(servers, exportPath)
	if err != nil {
		t.Fatalf("ExportServers failed: %v", err)
	}

	// Verify file exists and has correct permissions
	info, err := os.Stat(exportPath)
	if err != nil {
		t.Fatalf("Export file not found: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("File permissions = %o, want 0600", info.Mode().Perm())
	}

	// Verify it's valid JSON
	data, _ := os.ReadFile(exportPath)
	var parsed []map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Export JSON is invalid: %v", err)
	}
	if len(parsed) != 2 {
		t.Errorf("Exported %d servers, want 2", len(parsed))
	}
	if parsed[0]["name"] != "Server A" {
		t.Errorf("First server name = %v, want %q", parsed[0]["name"], "Server A")
	}
}

func TestImportServers(t *testing.T) {
	tmpDir := t.TempDir()
	importPath := filepath.Join(tmpDir, "import.json")

	importJSON := `[
		{
			"name": "Imported Server",
			"host": "192.168.1.50",
			"port": 22,
			"username": "admin",
			"password": "secret",
			"sftp_port": 22
		},
		{
			"name": "PEM Server",
			"host": "cloud.example.com",
			"port": 2222,
			"username": "ec2-user",
			"pem_key": "-----BEGIN RSA PRIVATE KEY-----\nMIIE...\n-----END RSA PRIVATE KEY-----"
		}
	]`
	os.WriteFile(importPath, []byte(importJSON), 0600)

	servers, err := ImportServers(importPath, 10)
	if err != nil {
		t.Fatalf("ImportServers failed: %v", err)
	}

	if len(servers) != 2 {
		t.Fatalf("Imported %d servers, want 2", len(servers))
	}

	// Check first server
	if servers[0].ID != 10 {
		t.Errorf("Server[0].ID = %d, want 10", servers[0].ID)
	}
	if servers[0].Name != "Imported Server" {
		t.Errorf("Server[0].Name = %q, want %q", servers[0].Name, "Imported Server")
	}
	if servers[0].Host != "192.168.1.50" {
		t.Errorf("Server[0].Host = %q, want %q", servers[0].Host, "192.168.1.50")
	}
	if servers[0].Password != "secret" {
		t.Errorf("Server[0].Password = %q, want %q", servers[0].Password, "secret")
	}

	// Check second server
	if servers[1].ID != 11 {
		t.Errorf("Server[1].ID = %d, want 11", servers[1].ID)
	}
	if servers[1].Port != 2222 {
		t.Errorf("Server[1].Port = %d, want 2222", servers[1].Port)
	}
	if servers[1].PemKey == "" {
		t.Error("Server[1].PemKey should not be empty")
	}
}

func TestImportServers_DefaultPort(t *testing.T) {
	tmpDir := t.TempDir()
	importPath := filepath.Join(tmpDir, "import.json")

	// No port specified — should default to 22
	importJSON := `[{"name": "No Port", "host": "10.0.0.1", "username": "root"}]`
	os.WriteFile(importPath, []byte(importJSON), 0600)

	servers, err := ImportServers(importPath, 1)
	if err != nil {
		t.Fatalf("ImportServers failed: %v", err)
	}
	if servers[0].Port != 22 {
		t.Errorf("Port = %d, want 22 (default)", servers[0].Port)
	}
}

func TestImportServers_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	importPath := filepath.Join(tmpDir, "bad.json")

	os.WriteFile(importPath, []byte("not json"), 0600)

	_, err := ImportServers(importPath, 1)
	if err == nil {
		t.Error("ImportServers should fail on invalid JSON")
	}
}

func TestImportServers_NonexistentFile(t *testing.T) {
	_, err := ImportServers("/nonexistent/path/import.json", 1)
	if err == nil {
		t.Error("ImportServers should fail on nonexistent file")
	}
}

func TestExportImportRoundTrip(t *testing.T) {
	tmpDir := t.TempDir()
	exportPath := filepath.Join(tmpDir, "roundtrip.json")

	original := []models.Server{
		{ID: 1, Name: "Roundtrip Server", Host: "10.0.0.1", Port: 443, Username: "admin", Password: "pw", SFTPPort: 443},
	}

	// Export
	err := ExportServers(original, exportPath)
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	// Import
	imported, err := ImportServers(exportPath, 100)
	if err != nil {
		t.Fatalf("Import failed: %v", err)
	}

	if len(imported) != 1 {
		t.Fatalf("Imported %d servers, want 1", len(imported))
	}

	// Verify data integrity (ID will differ since import assigns new IDs)
	if imported[0].Name != "Roundtrip Server" {
		t.Errorf("Name = %q, want %q", imported[0].Name, "Roundtrip Server")
	}
	if imported[0].Host != "10.0.0.1" {
		t.Errorf("Host = %q, want %q", imported[0].Host, "10.0.0.1")
	}
	if imported[0].Port != 443 {
		t.Errorf("Port = %d, want 443", imported[0].Port)
	}
	if imported[0].Username != "admin" {
		t.Errorf("Username = %q, want %q", imported[0].Username, "admin")
	}
	if imported[0].Password != "pw" {
		t.Errorf("Password = %q, want %q", imported[0].Password, "pw")
	}
}
