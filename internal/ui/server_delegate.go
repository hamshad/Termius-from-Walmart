package ui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"termius-from-walmart/internal/models"
)

// ServerDelegate is a custom list delegate for rendering server items
// with rich styling: icons, connection info, and auth badges.
type ServerDelegate struct{}

func NewServerDelegate() ServerDelegate {
	return ServerDelegate{}
}

func (d ServerDelegate) Height() int  { return 3 }
func (d ServerDelegate) Spacing() int { return 0 }
func (d ServerDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}

func (d ServerDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	server, ok := listItem.(models.Server)
	if !ok {
		return
	}

	isSelected := index == m.Index()

	// Icon
	icon := "  "
	if isSelected {
		icon = lipgloss.NewStyle().Foreground(ColorAccent).Bold(true).Render("> ")
	}

	// Server name
	nameStyle := lipgloss.NewStyle().Foreground(ColorText)
	if isSelected {
		nameStyle = lipgloss.NewStyle().Foreground(ColorAccent).Bold(true)
	}
	name := nameStyle.Render(server.Name)

	// Auth badge
	var authBadge string
	if server.PemKey != "" {
		authBadge = lipgloss.NewStyle().
			Foreground(ColorSuccess).
			Render(" [key]")
	} else if server.Password != "" {
		authBadge = lipgloss.NewStyle().
			Foreground(ColorHighlight).
			Render(" [pass]")
	} else {
		authBadge = lipgloss.NewStyle().
			Foreground(ColorTextMuted).
			Render(" [agent]")
	}

	// Connection string
	connStyle := lipgloss.NewStyle().Foreground(ColorTextDim)
	if isSelected {
		connStyle = lipgloss.NewStyle().Foreground(ColorTextMuted)
	}
	connStr := connStyle.Render(fmt.Sprintf("    %s@%s:%d", server.Username, server.Host, server.Port))

	// SFTP port hint (if different from SSH port)
	sftpHint := ""
	if server.SFTPPort != 0 && server.SFTPPort != server.Port {
		sftpHint = connStyle.Render(fmt.Sprintf("  sftp:%d", server.SFTPPort))
	}

	// Bottom separator
	sepStyle := lipgloss.NewStyle().Foreground(ColorBorder)
	width := m.Width()
	if width < 20 {
		width = 60
	}
	sep := sepStyle.Render(strings.Repeat("─", width-4))

	fmt.Fprintf(w, "%s%s%s\n%s%s\n%s", icon, name, authBadge, connStr, sftpHint, sep)
}
