package ui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"termius-from-walmart/internal/models"
	"termius-from-walmart/internal/storage"
)

// ViewState tracks which screen the TUI is currently showing.
type ViewState int

const (
	ListView ViewState = iota
	AddView
	EditView
	MenuView
	PemEditView
	FilePickerView
	SFTPView
	ConfirmView
)

// Model is the top-level Bubble Tea model for the entire application.
type Model struct {
	State       ViewState
	List        list.Model
	Config      *models.Config
	ConfigPath  string
	Inputs      []textinput.Model
	FocusIndex  int
	EditingID   int
	Message     string
	MenuOptions []string
	MenuCursor  int
	PemBuffer   string

	// File picker
	FilePickerList       list.Model
	FilePickerMode       string // "import" or "export"
	FilePickerPath       string
	FilePickerInput      textinput.Model
	FilePickerPrompt     bool
	FilePickerShowHidden bool

	// SFTP split-screen
	SelectedServer   *models.Server
	SFTPManager      *SFTPManager
	LocalFileList    list.Model
	RemoteFileList   list.Model
	LocalPath        string
	RemotePath       string
	FocusPane        string // "local" or "remote"
	TransferProgress int
	IsTransferring   bool
	TransferMessage  string

	// Confirmation dialog
	ConfirmTitle   string
	ConfirmMessage string
	ConfirmAction  func()
	ConfirmCursor  int // 0 = No, 1 = Yes
}

// NewModel creates and returns a fully initialized Model.
func NewModel() Model {
	configPath := storage.DefaultConfigPath()
	config := storage.LoadConfig(configPath)

	items := make([]list.Item, len(config.Servers))
	for i, server := range config.Servers {
		items[i] = server
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "SSH Connection Manager"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)

	return Model{
		State:       ListView,
		List:        l,
		Config:      config,
		ConfigPath:  configPath,
		MenuOptions: []string{"Import Servers", "Export Servers", "Back to List"},
		MenuCursor:  0,
		FilePickerList: func() list.Model {
			fl := list.New([]list.Item{}, FileDelegate{}, 0, 0)
			fl.SetShowStatusBar(false)
			fl.SetFilteringEnabled(false)
			return fl
		}(),
		FilePickerShowHidden: false,
		LocalFileList: func() list.Model {
			fl := list.New([]list.Item{}, FileDelegate{}, 0, 0)
			fl.SetShowStatusBar(false)
			fl.SetFilteringEnabled(false)
			return fl
		}(),
		RemoteFileList: func() list.Model {
			fl := list.New([]list.Item{}, FileDelegate{}, 0, 0)
			fl.SetShowStatusBar(false)
			fl.SetFilteringEnabled(false)
			return fl
		}(),
		LocalPath:        os.Getenv("HOME"),
		RemotePath:       "/",
		FocusPane:        "local",
		TransferProgress: 0,
		IsTransferring:   false,
	}
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := lipgloss.NewStyle().GetFrameSize()
		m.List.SetSize(msg.Width-h, msg.Height-v)
		m.FilePickerList.SetSize(msg.Width-h, msg.Height-v)
		halfWidth := (msg.Width - h - 2) / 2
		m.LocalFileList.SetSize(halfWidth, msg.Height-v-8)
		m.RemoteFileList.SetSize(halfWidth, msg.Height-v-8)
		return m, nil

	case tea.KeyMsg:
		switch m.State {
		case ListView:
			return m.updateListView(msg)
		case AddView, EditView:
			return m.updateFormView(msg)
		case FilePickerView:
			return m.updateFilePickerView(msg)
		case MenuView:
			return m.updateMenuView(msg)
		case PemEditView:
			return m.updatePemEditView(msg)
		case SFTPView:
			return m.updateSFTPView(msg)
		case ConfirmView:
			return m.updateConfirmView(msg)
		}
	}

	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	return m, cmd
}

// View implements tea.Model.
func (m Model) View() string {
	switch m.State {
	case ListView:
		return m.viewList()
	case AddView:
		return m.viewForm("Add New Server")
	case EditView:
		return m.viewForm("Edit Server")
	case MenuView:
		return m.viewMenu()
	case FilePickerView:
		return m.viewFilePicker()
	case PemEditView:
		return m.viewPemEdit()
	case SFTPView:
		return m.viewSFTP()
	case ConfirmView:
		return m.viewConfirm()
	}
	return ""
}

// RefreshList updates the bubbletea list from the current config.
func (m *Model) RefreshList() {
	items := make([]list.Item, len(m.Config.Servers))
	for i, server := range m.Config.Servers {
		items[i] = server
	}
	m.List.SetItems(items)
}

// SaveConfig persists the current config to disk.
func (m *Model) SaveConfig() error {
	return storage.SaveConfig(m.ConfigPath, m.Config)
}

// ShowConfirm switches to the confirmation dialog.
func (m *Model) ShowConfirm(title, message string, action func()) {
	m.ConfirmTitle = title
	m.ConfirmMessage = message
	m.ConfirmAction = action
	m.ConfirmCursor = 0 // default to "No"
	m.State = ConfirmView
}

// InitInputs creates fresh text input fields for the add/edit form.
func (m *Model) InitInputs() {
	m.Inputs = make([]textinput.Model, 7)

	m.Inputs[0] = textinput.New()
	m.Inputs[0].Placeholder = "Server Name"
	m.Inputs[0].Focus()
	m.Inputs[0].CharLimit = 50
	m.Inputs[0].Width = 40
	m.Inputs[0].Prompt = "Name: "

	m.Inputs[1] = textinput.New()
	m.Inputs[1].Placeholder = "192.168.1.1 or example.com"
	m.Inputs[1].CharLimit = 100
	m.Inputs[1].Width = 40
	m.Inputs[1].Prompt = "Host: "

	m.Inputs[2] = textinput.New()
	m.Inputs[2].Placeholder = "22"
	m.Inputs[2].CharLimit = 5
	m.Inputs[2].Width = 40
	m.Inputs[2].Prompt = "Port: "

	m.Inputs[3] = textinput.New()
	m.Inputs[3].Placeholder = "root"
	m.Inputs[3].CharLimit = 50
	m.Inputs[3].Width = 40
	m.Inputs[3].Prompt = "User: "

	m.Inputs[4] = textinput.New()
	m.Inputs[4].Placeholder = "password (optional, leave empty if using PEM)"
	m.Inputs[4].CharLimit = 100
	m.Inputs[4].Width = 40
	m.Inputs[4].Prompt = "Pass: "
	m.Inputs[4].EchoMode = textinput.EchoPassword
	m.Inputs[4].EchoCharacter = '*'

	m.Inputs[5] = textinput.New()
	m.Inputs[5].Placeholder = "Paste PEM key (optional, press 'p' to edit)"
	m.Inputs[5].CharLimit = 10000
	m.Inputs[5].Width = 40
	m.Inputs[5].Prompt = "PEM:  "

	m.Inputs[6] = textinput.New()
	m.Inputs[6].Placeholder = "22 (same as SSH port)"
	m.Inputs[6].CharLimit = 5
	m.Inputs[6].Width = 40
	m.Inputs[6].Prompt = "SFTP: "

	m.FocusIndex = 0
}

// PopulateInputsForEdit fills the form inputs with data from an existing server.
func (m *Model) PopulateInputsForEdit(server models.Server) {
	m.Inputs[0].SetValue(server.Name)
	m.Inputs[1].SetValue(server.Host)
	m.Inputs[2].SetValue(fmt.Sprintf("%d", server.Port))
	m.Inputs[3].SetValue(server.Username)
	m.Inputs[4].SetValue(server.Password)
	m.Inputs[5].SetValue(server.PemKey)
	if server.SFTPPort == 0 {
		m.Inputs[6].SetValue(fmt.Sprintf("%d", server.Port))
	} else {
		m.Inputs[6].SetValue(fmt.Sprintf("%d", server.SFTPPort))
	}
}
