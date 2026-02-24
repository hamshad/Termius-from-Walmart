package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"termius-from-walmart/internal/storage"
)

// updateMenuView handles key events in the import/export menu.
func (m Model) updateMenuView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		m.State = ListView
		return m, nil

	case "up", "k":
		if m.MenuCursor > 0 {
			m.MenuCursor--
		}

	case "down", "j":
		if m.MenuCursor < len(m.MenuOptions)-1 {
			m.MenuCursor++
		}

	case "enter":
		switch m.MenuCursor {
		case 0: // Import
			m.State = FilePickerView
			m.FilePickerMode = "import"
			m.FilePickerPath = os.Getenv("HOME")
			m.FilePickerPrompt = false
			m.loadFileList()
		case 1: // Export
			m.State = FilePickerView
			m.FilePickerMode = "export"
			m.FilePickerPath = os.Getenv("HOME")
			m.FilePickerPrompt = false
			ti := textinput.New()
			ti.Placeholder = "ssh-servers-export.json"
			ti.Width = 40
			ti.Prompt = "Filename: "
			m.FilePickerInput = ti
			m.loadFileList()
		case 2: // Back
			m.State = ListView
		}
		return m, nil

	case "esc":
		m.State = ListView
		return m, nil
	}

	return m, nil
}

// updateFilePickerView handles key events in the file picker.
func (m Model) updateFilePickerView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.FilePickerPrompt {
		var cmd tea.Cmd
		m.FilePickerInput, cmd = m.FilePickerInput.Update(msg)
		if msg.String() == "enter" {
			filename := strings.TrimSpace(m.FilePickerInput.Value())
			if filename != "" {
				full := filepath.Join(m.FilePickerPath, filename)
				m.exportServersToPath(full)
				m.FilePickerPrompt = false
				m.State = ListView
			} else {
				m.Message = "Filename cannot be empty"
			}
		} else if msg.String() == "esc" {
			m.FilePickerPrompt = false
		}
		return m, cmd
	}

	switch msg.String() {
	case "ctrl+c", "q", "esc":
		m.State = ListView
		return m, nil
	case ".":
		m.FilePickerShowHidden = !m.FilePickerShowHidden
		m.loadFileList()
		return m, nil
	case "enter":
		sel := m.FilePickerList.SelectedItem()
		if sel == nil {
			return m, nil
		}
		name := sel.FilterValue()
		if name == "../" {
			parent := filepath.Dir(m.FilePickerPath)
			m.FilePickerPath = parent
			m.loadFileList()
			return m, nil
		}
		if strings.HasSuffix(name, "/") {
			dirName := strings.TrimSuffix(name, "/")
			m.FilePickerPath = filepath.Join(m.FilePickerPath, dirName)
			m.loadFileList()
			return m, nil
		}
		full := filepath.Join(m.FilePickerPath, name)
		if m.FilePickerMode == "import" {
			m.importServersFromPath(full)
			m.State = ListView
			return m, nil
		}
		if m.FilePickerMode == "export" {
			m.exportServersToPath(full)
			m.State = ListView
			return m, nil
		}

	case "x":
		if m.FilePickerMode == "export" {
			full := filepath.Join(m.FilePickerPath, "ssh-servers-export.json")
			m.exportServersToPath(full)
			m.State = ListView
			return m, nil
		}
	case "n":
		if m.FilePickerMode == "export" {
			m.FilePickerPrompt = true
			m.FilePickerInput.SetValue("")
			m.FilePickerInput.Focus()
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.FilePickerList, cmd = m.FilePickerList.Update(msg)
	return m, cmd
}

func (m *Model) loadFileList() {
	entries, err := os.ReadDir(m.FilePickerPath)
	if err != nil {
		m.FilePickerList.SetItems([]list.Item{})
		m.Message = fmt.Sprintf("Unable to read %s: %v", m.FilePickerPath, err)
		return
	}

	items := make([]list.Item, 0, len(entries)+1)
	if parent := filepath.Dir(m.FilePickerPath); parent != m.FilePickerPath {
		items = append(items, FileItem("../"))
	}

	for _, e := range entries {
		name := e.Name()
		if !m.FilePickerShowHidden && strings.HasPrefix(name, ".") {
			continue
		}
		if e.IsDir() {
			name = name + "/"
		}
		items = append(items, FileItem(name))
	}

	m.FilePickerList.SetItems(items)
	m.FilePickerList.Title = fmt.Sprintf("Select file (%s)", m.FilePickerMode)
}

func (m *Model) importServersFromPath(importPath string) {
	imported, err := storage.ImportServers(importPath, m.Config.NextID)
	if err != nil {
		m.Message = fmt.Sprintf("Import failed: %v", err)
		return
	}

	m.Config.Servers = append(m.Config.Servers, imported...)
	m.Config.NextID += len(imported)

	if err := m.SaveConfig(); err != nil {
		m.Message = fmt.Sprintf("Import failed: %v", err)
		return
	}

	m.RefreshList()
	m.Message = fmt.Sprintf("Imported %d servers from %s", len(imported), importPath)
}

func (m *Model) exportServersToPath(exportPath string) {
	if err := storage.ExportServers(m.Config.Servers, exportPath); err != nil {
		m.Message = fmt.Sprintf("Export failed: %v", err)
		return
	}
	m.Message = fmt.Sprintf("Exported %d servers to %s", len(m.Config.Servers), exportPath)
}
