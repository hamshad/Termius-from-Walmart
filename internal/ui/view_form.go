package ui

import (
	"fmt"
	"strconv"
	"strings"

	"termius-from-walmart/internal/models"
	"termius-from-walmart/internal/ssh"

	tea "github.com/charmbracelet/bubbletea"
)

// updateFormView handles key events in the add/edit server form.
func (m Model) updateFormView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit

	case "esc":
		m.State = ListView
		m.Message = ""
		return m, nil

	case "p":
		if m.FocusIndex == 5 {
			m.PemBuffer = m.Inputs[5].Value()
			prevState := m.State
			m.State = PemEditView
			if prevState == AddView {
				m.EditingID = -1
			}
			return m, nil
		}

	case "tab", "shift+tab", "up", "down":
		s := msg.String()
		if s == "up" || s == "shift+tab" {
			m.FocusIndex--
		} else {
			m.FocusIndex++
		}

		if m.FocusIndex > len(m.Inputs) {
			m.FocusIndex = 0
		} else if m.FocusIndex < 0 {
			m.FocusIndex = len(m.Inputs)
		}

		for i := 0; i < len(m.Inputs); i++ {
			if i == m.FocusIndex {
				m.Inputs[i].Focus()
			} else {
				m.Inputs[i].Blur()
			}
		}

		return m, nil

	case "enter":
		if m.FocusIndex == len(m.Inputs) {
			if m.validateAndSave() {
				m.State = ListView
			}
			return m, nil
		}
	}

	cmd := m.updateInputs(msg)
	return m, cmd
}

// updatePemEditView handles key events in the PEM key multiline editor.
func (m Model) updatePemEditView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit

	case "ctrl+s":
		m.Inputs[5].SetValue(m.PemBuffer)
		if m.EditingID == -1 {
			m.State = AddView
		} else {
			m.State = EditView
		}
		m.Message = "PEM key saved"
		return m, nil

	case "esc":
		if m.EditingID == -1 {
			m.State = AddView
		} else {
			m.State = EditView
		}
		return m, nil

	default:
		key := msg.String()
		if key == "backspace" {
			if len(m.PemBuffer) > 0 {
				m.PemBuffer = m.PemBuffer[:len(m.PemBuffer)-1]
			}
		} else if key == "enter" {
			m.PemBuffer += "\n"
		} else if len(key) == 1 {
			m.PemBuffer += key
		}
	}

	return m, nil
}

func (m *Model) updateInputs(msg tea.KeyMsg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.Inputs))
	for i := range m.Inputs {
		m.Inputs[i], cmds[i] = m.Inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func (m *Model) validateAndSave() bool {
	name := strings.TrimSpace(m.Inputs[0].Value())
	host := strings.TrimSpace(m.Inputs[1].Value())
	portStr := strings.TrimSpace(m.Inputs[2].Value())
	username := strings.TrimSpace(m.Inputs[3].Value())
	password := m.Inputs[4].Value()
	pemKey := m.Inputs[5].Value()
	sftpPortStr := strings.TrimSpace(m.Inputs[6].Value())

	if name == "" {
		m.Message = "Error: Name is required"
		return false
	}
	if host == "" {
		m.Message = "Error: Host is required"
		return false
	}
	if username == "" {
		m.Message = "Error: Username is required"
		return false
	}
	if password != "" && pemKey != "" {
		m.Message = "Error: Use either password OR PEM key, not both"
		return false
	}

	port := 22
	if portStr != "" {
		var err error
		port, err = strconv.Atoi(portStr)
		if err != nil || port < 1 || port > 65535 {
			m.Message = "Error: Invalid port number"
			return false
		}
	}

	sftpPort := port
	if sftpPortStr != "" {
		var err error
		sftpPort, err = strconv.Atoi(sftpPortStr)
		if err != nil || sftpPort < 1 || sftpPort > 65535 {
			m.Message = "Error: Invalid SFTP port number"
			return false
		}
	}

	if pemKey != "" {
		normalized := ssh.NormalizePemKey(pemKey)
		if !strings.Contains(normalized, "BEGIN") || !strings.Contains(normalized, "PRIVATE KEY") {
			m.Message = "Error: Invalid PEM key format"
			return false
		}
	}

	if m.State == AddView {
		server := models.Server{
			ID:       m.Config.NextID,
			Name:     name,
			Host:     host,
			Port:     port,
			Username: username,
			Password: password,
			PemKey:   pemKey,
			SFTPPort: sftpPort,
		}
		m.Config.Servers = append(m.Config.Servers, server)
		m.Config.NextID++
		m.Message = fmt.Sprintf("Added server: %s", name)
	} else if m.State == EditView {
		for i, server := range m.Config.Servers {
			if server.ID == m.EditingID {
				m.Config.Servers[i].Name = name
				m.Config.Servers[i].Host = host
				m.Config.Servers[i].Port = port
				m.Config.Servers[i].Username = username
				m.Config.Servers[i].Password = password
				m.Config.Servers[i].PemKey = pemKey
				m.Config.Servers[i].SFTPPort = sftpPort
				m.Message = fmt.Sprintf("Updated server: %s", name)
				break
			}
		}
	}

	if err := m.SaveConfig(); err != nil {
		m.Message = fmt.Sprintf("Error saving: %v", err)
		return false
	}

	m.RefreshList()
	return true
}
