package ui

import (
	"fmt"

	"termius-from-walmart/internal/models"
	"termius-from-walmart/internal/ssh"

	tea "github.com/charmbracelet/bubbletea"
)

// updateListView handles key events in the main server list view.
func (m Model) updateListView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "a":
		m.State = AddView
		m.InitInputs()
		m.Message = ""
		return m, nil

	case "e":
		if len(m.Config.Servers) > 0 {
			selected := m.List.SelectedItem()
			if server, ok := selected.(models.Server); ok {
				m.State = EditView
				m.EditingID = server.ID
				m.InitInputs()
				m.PopulateInputsForEdit(server)
				m.Message = ""
				return m, nil
			}
		}

	case "d":
		if len(m.Config.Servers) > 0 {
			selected := m.List.SelectedItem()
			if server, ok := selected.(models.Server); ok {
				serverToDelete := server
				m.ShowConfirm(
					"Delete Server",
					fmt.Sprintf("Are you sure you want to delete '%s' (%s@%s)?", server.Name, server.Username, server.Host),
					func() {
						newServers := []models.Server{}
						for _, s := range m.Config.Servers {
							if s.ID != serverToDelete.ID {
								newServers = append(newServers, s)
							}
						}
						m.Config.Servers = newServers
						m.SaveConfig()
						m.RefreshList()
						m.Message = fmt.Sprintf("Deleted server: %s", serverToDelete.Name)
					},
				)
				return m, nil
			}
		}

	case "enter":
		if len(m.Config.Servers) > 0 {
			selected := m.List.SelectedItem()
			if server, ok := selected.(models.Server); ok {
				m.Message = fmt.Sprintf("Connecting to %s@%s:%d...", server.Username, server.Host, server.Port)
				return m, tea.Sequence(
					tea.ExecProcess(ssh.BuildSSHCommand(server), func(err error) tea.Msg {
						return err
					}),
				)
			}
		}

	case "m":
		m.State = MenuView
		m.MenuCursor = 0
		m.Message = ""
		return m, nil

	case "s":
		if len(m.Config.Servers) > 0 {
			selected := m.List.SelectedItem()
			if server, ok := selected.(models.Server); ok {
				m.SelectedServer = &server
				m.ConnectingName = server.Name
				m.SpinnerFrame = 0
				m.State = ConnectingView
				m.Message = ""
				return m, tea.Batch(connectSFTPCmd(&server), spinnerTick())
			}
		}
	}

	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	return m, cmd
}
