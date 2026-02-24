package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// updateSFTPView handles key events in the SFTP split-screen view.
func (m Model) updateSFTPView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q", "esc":
		if m.SFTPManager != nil {
			m.SFTPManager.Close()
		}
		m.SelectedServer = nil
		m.SFTPManager = nil
		m.State = ListView
		m.Message = ""
		return m, nil

	case "tab":
		if m.FocusPane == "local" {
			m.FocusPane = "remote"
		} else {
			m.FocusPane = "local"
		}
		return m, nil

	case "enter":
		if m.FocusPane == "local" {
			sel := m.LocalFileList.SelectedItem()
			if sel == nil {
				return m, nil
			}
			m.navigateLocalDir(sel.FilterValue())
		} else {
			sel := m.RemoteFileList.SelectedItem()
			if sel == nil {
				return m, nil
			}
			m.navigateRemoteDir(sel.FilterValue())
		}
		return m, nil

	case "c":
		if m.IsTransferring {
			return m, nil // don't start another transfer while one is in progress
		}
		cmd := m.performCopy()
		return m, cmd

	case "d":
		if m.FocusPane == "local" {
			sel := m.LocalFileList.SelectedItem()
			if sel != nil {
				fileName := sel.FilterValue()
				if !strings.HasSuffix(fileName, "/") {
					m.ShowConfirm(
						"Delete Local File",
						fmt.Sprintf("Delete '%s'?", fileName),
						func() {
							filePath := filepath.Join(m.LocalPath, fileName)
							if err := os.Remove(filePath); err != nil {
								m.Message = fmt.Sprintf("Error deleting: %v", err)
							} else {
								m.Message = fmt.Sprintf("Deleted %s", fileName)
								m.loadLocalFiles(m.LocalPath)
							}
						},
					)
				} else {
					m.Message = "Cannot delete directories"
				}
			}
		} else {
			sel := m.RemoteFileList.SelectedItem()
			if sel != nil {
				fileName := sel.FilterValue()
				if !strings.HasSuffix(fileName, "/") {
					m.ShowConfirm(
						"Delete Remote File",
						fmt.Sprintf("Delete '%s' from remote server?", fileName),
						func() {
							if m.SFTPManager == nil {
								m.Message = "Error: SFTP connection lost"
								return
							}
							filePath := filepath.Join(m.RemotePath, fileName)
							if err := m.SFTPManager.DeleteFile(filePath); err != nil {
								m.Message = fmt.Sprintf("Error deleting: %v", err)
							} else {
								m.Message = fmt.Sprintf("Deleted %s", fileName)
								m.loadRemoteFiles(m.RemotePath)
							}
						},
					)
				} else {
					m.Message = "Cannot delete directories"
				}
			}
		}
		return m, nil

	case "r":
		m.Message = "Rename not yet implemented"
		return m, nil
	}

	var cmd tea.Cmd
	if m.FocusPane == "local" {
		m.LocalFileList, cmd = m.LocalFileList.Update(msg)
	} else {
		m.RemoteFileList, cmd = m.RemoteFileList.Update(msg)
	}
	return m, cmd
}

func (m *Model) loadLocalFiles(path string) {
	entries, err := os.ReadDir(path)
	if err != nil {
		m.Message = fmt.Sprintf("Error reading local directory: %v", err)
		m.LocalFileList.SetItems([]list.Item{})
		return
	}

	items := make([]list.Item, 0, len(entries)+1)
	if path != "/" {
		items = append(items, FileItem{Name: "..", IsDir: true})
	}

	for _, f := range entries {
		items = append(items, NewFileItemFromDirEntry(f))
	}

	m.LocalFileList.SetItems(items)
	m.LocalPath = path
}

func (m *Model) loadRemoteFiles(path string) {
	if m.SFTPManager == nil {
		m.Message = "Error: SFTP connection lost"
		return
	}

	files, err := m.SFTPManager.ListFiles(path)
	if err != nil {
		m.Message = fmt.Sprintf("Error listing remote files: %v", err)
		m.RemoteFileList.SetItems([]list.Item{})
		return
	}

	items := make([]list.Item, 0, len(files)+1)
	if path != "/" {
		items = append(items, FileItem{Name: "..", IsDir: true})
	}

	for _, f := range files {
		items = append(items, NewFileItemFromInfo(f))
	}

	m.RemoteFileList.SetItems(items)
	m.RemotePath = path
}

func (m *Model) navigateLocalDir(fileName string) {
	if fileName == "../" {
		m.LocalPath = filepath.Dir(m.LocalPath)
		m.loadLocalFiles(m.LocalPath)
		return
	}
	if strings.HasSuffix(fileName, "/") {
		dirName := strings.TrimSuffix(fileName, "/")
		m.LocalPath = filepath.Join(m.LocalPath, dirName)
		m.loadLocalFiles(m.LocalPath)
	}
}

func (m *Model) navigateRemoteDir(fileName string) {
	if fileName == "../" {
		parent := filepath.Dir(m.RemotePath)
		if parent == "." {
			parent = "/"
		}
		m.RemotePath = parent
		m.loadRemoteFiles(m.RemotePath)
		return
	}
	if strings.HasSuffix(fileName, "/") {
		dirName := strings.TrimSuffix(fileName, "/")
		var newPath string
		if m.RemotePath == "/" {
			newPath = "/" + dirName
		} else {
			newPath = filepath.Join(m.RemotePath, dirName)
		}
		m.RemotePath = newPath
		m.loadRemoteFiles(m.RemotePath)
	}
}

func (m *Model) performCopy() tea.Cmd {
	var srcFile, dstFile string
	var isLocalToRemote bool

	if m.FocusPane == "local" {
		sel := m.LocalFileList.SelectedItem()
		if sel == nil {
			m.Message = "No file selected"
			return nil
		}
		fileName := sel.FilterValue()
		if strings.HasSuffix(fileName, "/") {
			m.Message = "Cannot copy directories"
			return nil
		}
		srcFile = filepath.Join(m.LocalPath, fileName)
		dstFile = filepath.Join(m.RemotePath, fileName)
		isLocalToRemote = true
	} else {
		sel := m.RemoteFileList.SelectedItem()
		if sel == nil {
			m.Message = "No file selected"
			return nil
		}
		fileName := sel.FilterValue()
		if strings.HasSuffix(fileName, "/") {
			m.Message = "Cannot copy directories"
			return nil
		}
		srcFile = filepath.Join(m.RemotePath, fileName)
		dstFile = filepath.Join(m.LocalPath, fileName)
		isLocalToRemote = false
	}

	if m.SFTPManager == nil {
		m.Message = "Error: SFTP connection lost"
		return nil
	}

	m.IsTransferring = true
	m.TransferMessage = fmt.Sprintf("Copying %s...", filepath.Base(srcFile))
	m.TransferProgress = 0

	return transferFileCmd(m.SFTPManager, srcFile, dstFile, isLocalToRemote)
}
