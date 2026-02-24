package models

import "fmt"

// Server represents a saved SSH server configuration.
type Server struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	PemKey   string `json:"pem_key"`
	SFTPPort int    `json:"sftp_port"`
}

// FilterValue implements list.Item for bubbletea list filtering.
func (s Server) FilterValue() string { return s.Name }

// Title implements list.Item for bubbletea list display.
func (s Server) Title() string { return s.Name }

// Description implements list.Item for bubbletea list display.
func (s Server) Description() string {
	return fmt.Sprintf("%s@%s:%d", s.Username, s.Host, s.Port)
}

// Config holds all servers and the auto-incrementing ID counter.
type Config struct {
	Servers []Server `json:"servers"`
	NextID  int      `json:"next_id"`
}

// NewConfig returns a default empty configuration.
func NewConfig() *Config {
	return &Config{
		Servers: []Server{},
		NextID:  1,
	}
}
