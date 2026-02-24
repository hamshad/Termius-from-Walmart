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

	fn := FileItemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return FileSelectedStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(name))
}
