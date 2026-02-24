package ui

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// FileItem represents a file or directory entry in a file picker / SFTP list.
type FileItem struct {
	Name    string
	Size    int64
	Mode    os.FileMode
	ModTime time.Time
	IsDir   bool
}

func (f FileItem) FilterValue() string {
	if f.IsDir {
		return f.Name + "/"
	}
	return f.Name
}
func (f FileItem) Title() string       { return f.FilterValue() }
func (f FileItem) Description() string { return "" }

// NewFileItem creates a FileItem from a name string (legacy compat for file picker).
func NewFileItem(name string) FileItem {
	isDir := strings.HasSuffix(name, "/")
	cleanName := strings.TrimSuffix(name, "/")
	return FileItem{Name: cleanName, IsDir: isDir}
}

// NewFileItemFromInfo creates a FileItem from os.FileInfo (for SFTP remote files).
func NewFileItemFromInfo(info os.FileInfo) FileItem {
	return FileItem{
		Name:    info.Name(),
		Size:    info.Size(),
		Mode:    info.Mode(),
		ModTime: info.ModTime(),
		IsDir:   info.IsDir(),
	}
}

// NewFileItemFromDirEntry creates a FileItem from os.DirEntry (for local files).
func NewFileItemFromDirEntry(entry os.DirEntry) FileItem {
	fi := FileItem{
		Name:  entry.Name(),
		IsDir: entry.IsDir(),
	}
	if info, err := entry.Info(); err == nil {
		fi.Size = info.Size()
		fi.Mode = info.Mode()
		fi.ModTime = info.ModTime()
	}
	return fi
}

// FileDelegate is a compact list delegate for file items with metadata.
type FileDelegate struct {
	ShowMetadata bool // whether to show size/perms/date
}

func (d FileDelegate) Height() int {
	if d.ShowMetadata {
		return 2
	}
	return 1
}
func (d FileDelegate) Spacing() int                            { return 0 }
func (d FileDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d FileDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	fi, ok := listItem.(FileItem)
	if !ok {
		return
	}

	displayName := fi.Name
	if fi.IsDir {
		displayName = fi.Name + "/"
	}
	isSelected := index == m.Index()

	var rendered string
	switch {
	case fi.IsDir && isSelected:
		rendered = FileDirSelectedStyle.Render("> " + displayName)
	case fi.IsDir:
		rendered = FileDirStyle.Render("  " + displayName)
	case isSelected:
		rendered = FileSelectedStyle.Render("> " + displayName)
	default:
		rendered = FileItemStyle.Render("  " + displayName)
	}

	fmt.Fprint(w, rendered)

	// Show metadata on second line for non-directory files when enabled
	if d.ShowMetadata && !fi.IsDir && fi.Name != ".." {
		meta := formatFileMetadata(fi, m.Width())
		if isSelected {
			fmt.Fprint(w, "\n"+lipgloss.NewStyle().Foreground(ColorTextMuted).PaddingLeft(4).Render(meta))
		} else {
			fmt.Fprint(w, "\n"+lipgloss.NewStyle().Foreground(ColorTextMuted).PaddingLeft(4).Render(meta))
		}
	} else if d.ShowMetadata {
		// Still need to output the second line for consistent height
		fmt.Fprint(w, "\n")
	}
}

// formatFileMetadata formats file size, permissions, and mod time into a compact string.
func formatFileMetadata(fi FileItem, maxWidth int) string {
	size := humanizeBytes(fi.Size)
	perms := fi.Mode.Perm().String()
	modTime := fi.ModTime.Format("Jan 02 15:04")

	// Compact format
	meta := fmt.Sprintf("%s  %s  %s", size, perms, modTime)
	if maxWidth > 0 && len(meta) > maxWidth-6 {
		meta = fmt.Sprintf("%s  %s", size, modTime)
	}
	return meta
}

// humanizeBytes converts byte count to human-readable string.
func humanizeBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%dB", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%c", float64(b)/float64(div), "KMGTPE"[exp])
}
