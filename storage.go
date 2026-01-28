package main

import (
    "encoding/json"
    "io/ioutil"
    "os"
    "path/filepath"
)

// loadConfig reads the config from the given path. If the file does not exist
// or cannot be read, it returns a default empty config. The config directory
// will be created with 0700 permissions if missing.
func loadConfig(path string) *Config {
    config := &Config{
        Servers: []Server{},
        NextID:  1,
    }

    dir := filepath.Dir(path)
    if err := os.MkdirAll(dir, 0700); err != nil {
        return config
    }

    data, err := ioutil.ReadFile(path)
    if err != nil {
        return config
    }

    if err := json.Unmarshal(data, config); err != nil {
        return config
    }
    return config
}

// saveConfig writes the current config to disk with 0600 permissions.
func (m *model) saveConfig() error {
    data, err := json.MarshalIndent(m.config, "", "  ")
    if err != nil {
        return err
    }

    dir := filepath.Dir(m.configPath)
    if err := os.MkdirAll(dir, 0700); err != nil {
        return err
    }

    return ioutil.WriteFile(m.configPath, data, 0600)
}
