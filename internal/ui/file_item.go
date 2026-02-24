package ui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// FileItem represents a file or directory entry in a file picker / SFTP list.
type FileItem string

func (f FileItem) FilterValue() string { return string(f) }
func (f FileItem) Title() string       { return string(f) }
func (f FileItem) Description() string { return "" }

// FileDelegate is a compact single-line list delegate for file items.
type FileDelegate struct{}

func (d FileDelegate) Height() int                             { return 1 }
func (d FileDelegate) Spacing() int                            { return 0 }
func (d FileDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d FileDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	fi, ok := listItem.(FileItem)
	if !ok {
		return
	}
	name := string(fi)
	isDir := strings.HasSuffix(name, "/")
	isSelected := index == m.Index()

	var rendered string
	switch {
	case isDir && isSelected:
		rendered = FileDirSelectedStyle.Render("> " + name)
	case isDir:
		rendered = FileDirStyle.Render("  " + name)
	case isSelected:
		rendered = FileSelectedStyle.Render("> " + name)
	default:
		rendered = FileItemStyle.Render("  " + name)
	}

	fmt.Fprint(w, rendered)
}
