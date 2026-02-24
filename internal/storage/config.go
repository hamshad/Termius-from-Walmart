package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"termius-from-walmart/internal/models"
)

// DefaultConfigPath returns the default config file path under the user's home directory.
func DefaultConfigPath() string {
	return filepath.Join(os.Getenv("HOME"), ".termius-from-walmart", "config.json")
}

// LoadConfig reads the config from the given path. If the file does not exist
// or cannot be parsed, it returns a default empty config. The config directory
// is created with 0700 permissions if missing.
func LoadConfig(path string) *models.Config {
	config := models.NewConfig()

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return config
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return config
	}

	if err := json.Unmarshal(data, config); err != nil {
		return config
	}
	return config
}

// SaveConfig writes the config to disk with 0600 permissions.
func SaveConfig(path string, config *models.Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	return os.WriteFile(path, data, 0600)
}

// ExportServers exports the server list to a JSON file at the given path.
func ExportServers(servers []models.Server, exportPath string) error {
	exportData := make([]map[string]interface{}, len(servers))
	for i, server := range servers {
		exportData[i] = map[string]interface{}{
			"name":      server.Name,
			"host":      server.Host,
			"port":      server.Port,
			"username":  server.Username,
			"password":  server.Password,
			"pem_key":   server.PemKey,
			"sftp_port": server.SFTPPort,
		}
	}

	data, err := json.MarshalIndent(exportData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal export data: %w", err)
	}

	return os.WriteFile(exportPath, data, 0600)
}

// ImportServers reads a JSON file and returns a slice of servers parsed from it.
// Each server is assigned an ID starting from nextID.
func ImportServers(importPath string, nextID int) ([]models.Server, error) {
	data, err := os.ReadFile(importPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read import file: %w", err)
	}

	var importData []map[string]interface{}
	if err := json.Unmarshal(data, &importData); err != nil {
		return nil, fmt.Errorf("invalid JSON format: %w", err)
	}

	servers := make([]models.Server, 0, len(importData))
	for _, item := range importData {
		server := models.Server{
			ID:       nextID,
			Name:     fmt.Sprintf("%v", item["name"]),
			Host:     fmt.Sprintf("%v", item["host"]),
			Username: fmt.Sprintf("%v", item["username"]),
		}

		if port, ok := item["port"].(float64); ok {
			server.Port = int(port)
		} else {
			server.Port = 22
		}

		if password, ok := item["password"].(string); ok {
			server.Password = password
		}

		if pemKey, ok := item["pem_key"].(string); ok {
			server.PemKey = pemKey
		}

		if sftpPort, ok := item["sftp_port"].(float64); ok {
			server.SFTPPort = int(sftpPort)
		} else {
			server.SFTPPort = server.Port
		}

		servers = append(servers, server)
		nextID++
	}

	return servers, nil
}
