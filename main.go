package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Server represents a saved SSH server configuration
type Server struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	PemKey   string `json:"pem_key"` // PEM private key as text
}

// Config holds all servers and keychains
type Config struct {
	Servers []Server `json:"servers"`
	NextID  int      `json:"next_id"`
}

// Implement list.Item interface for Server
func (s Server) FilterValue() string { return s.Name }
func (s Server) Title() string       { return s.Name }
func (s Server) Description() string { return fmt.Sprintf("%s@%s:%d", s.Username, s.Host, s.Port) }

// View states
type viewState int

const (
	listView viewState = iota
	addView
	editView
	menuView
	pemEditView
	filePickerView
)

type model struct {
	state       viewState
	list        list.Model
	config      *Config
	configPath  string
	inputs      []textinput.Model
	focusIndex  int
	editingID   int
	message     string
	menuOptions []string
	menuCursor  int
	pemBuffer   string // Buffer for multiline PEM editing
	// File picker fields
	filePickerList       list.Model
	filePickerMode       string // "import" or "export"
	filePickerPath       string
	filePickerInput      textinput.Model
	filePickerPrompt     bool // whether we're typing a filename for export
	filePickerShowHidden bool // whether to show dotfiles
}

var (
	titleStyle = lipgloss.NewStyle().
			MarginLeft(2).
			Foreground(lipgloss.Color("62")).
			Bold(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginLeft(2)

	messageStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("170")).
			MarginLeft(2).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			MarginLeft(2).
			Bold(true)
)

var (
	fileItemStyle     = lipgloss.NewStyle().PaddingLeft(2)
	fileSelectedStyle = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("170"))
)

func initialModel() model {
	configPath := filepath.Join(os.Getenv("HOME"), ".termius-from-walmart", "config.json")
	config := loadConfig(configPath)

	items := make([]list.Item, len(config.Servers))
	for i, server := range config.Servers {
		items[i] = server
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "SSH Connection Manager"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)

	return model{
		state:       listView,
		list:        l,
		config:      config,
		configPath:  configPath,
		menuOptions: []string{"Import Servers", "Export Servers", "Back to List"},
		menuCursor:  0,
		// create file picker list with compact delegate
		filePickerList: func() list.Model {
			l := list.New([]list.Item{}, fileDelegate{}, 0, 0)
			l.SetShowStatusBar(false)
			l.SetFilteringEnabled(false)
			return l
		}(),
		filePickerShowHidden: false,
	}
}

func (m *model) initInputs() {
	m.inputs = make([]textinput.Model, 6)

	// Name
	m.inputs[0] = textinput.New()
	m.inputs[0].Placeholder = "Server Name"
	m.inputs[0].Focus()
	m.inputs[0].CharLimit = 50
	m.inputs[0].Width = 40
	m.inputs[0].Prompt = "Name: "

	// Host
	m.inputs[1] = textinput.New()
	m.inputs[1].Placeholder = "192.168.1.1 or example.com"
	m.inputs[1].CharLimit = 100
	m.inputs[1].Width = 40
	m.inputs[1].Prompt = "Host: "

	// Port
	m.inputs[2] = textinput.New()
	m.inputs[2].Placeholder = "22"
	m.inputs[2].CharLimit = 5
	m.inputs[2].Width = 40
	m.inputs[2].Prompt = "Port: "

	// Username
	m.inputs[3] = textinput.New()
	m.inputs[3].Placeholder = "root"
	m.inputs[3].CharLimit = 50
	m.inputs[3].Width = 40
	m.inputs[3].Prompt = "User: "

	// Password
	m.inputs[4] = textinput.New()
	m.inputs[4].Placeholder = "password (optional, leave empty if using PEM)"
	m.inputs[4].CharLimit = 100
	m.inputs[4].Width = 40
	m.inputs[4].Prompt = "Pass: "
	m.inputs[4].EchoMode = textinput.EchoPassword
	m.inputs[4].EchoCharacter = '•'

	// PEM Key
	m.inputs[5] = textinput.New()
	m.inputs[5].Placeholder = "Paste PEM key (optional, press 'p' to edit)"
	m.inputs[5].CharLimit = 10000
	m.inputs[5].Width = 40
	m.inputs[5].Prompt = "PEM:  "

	m.focusIndex = 0
}

func (m *model) populateInputsForEdit(server Server) {
	m.inputs[0].SetValue(server.Name)
	m.inputs[1].SetValue(server.Host)
	m.inputs[2].SetValue(strconv.Itoa(server.Port))
	m.inputs[3].SetValue(server.Username)
	m.inputs[4].SetValue(server.Password)
	m.inputs[5].SetValue(server.PemKey)
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := lipgloss.NewStyle().GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
		m.filePickerList.SetSize(msg.Width-h, msg.Height-v)
		return m, nil

	case tea.KeyMsg:
		switch m.state {
		case listView:
			return m.updateListView(msg)
		case addView, editView:
			return m.updateFormView(msg)
		case filePickerView:
			return m.updateFilePickerView(msg)
		case menuView:
			return m.updateMenuView(msg)
		case pemEditView:
			return m.updatePemEditView(msg)
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) updateListView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "a":
		m.state = addView
		m.initInputs()
		m.message = ""
		return m, nil

	case "e":
		if len(m.config.Servers) > 0 {
			selected := m.list.SelectedItem()
			if server, ok := selected.(Server); ok {
				m.state = editView
				m.editingID = server.ID
				m.initInputs()
				m.populateInputsForEdit(server)
				m.message = ""
				return m, nil
			}
		}

	case "d":
		if len(m.config.Servers) > 0 {
			selected := m.list.SelectedItem()
			if server, ok := selected.(Server); ok {
				m.deleteServer(server.ID)
				m.message = fmt.Sprintf("Deleted server: %s", server.Name)
				return m, nil
			}
		}

	case "enter":
		if len(m.config.Servers) > 0 {
			selected := m.list.SelectedItem()
			if server, ok := selected.(Server); ok {
				return m, tea.Sequence(
					tea.ExecProcess(m.connectSSH(server), func(err error) tea.Msg {
						return err
					}),
				)
			}
		}

	case "m":
		m.state = menuView
		m.menuCursor = 0
		m.message = ""
		return m, nil
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) updateMenuView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		m.state = listView
		return m, nil

	case "up", "k":
		if m.menuCursor > 0 {
			m.menuCursor--
		}

	case "down", "j":
		if m.menuCursor < len(m.menuOptions)-1 {
			m.menuCursor++
		}

	case "enter":
		switch m.menuCursor {
		case 0: // Import
			m.state = filePickerView
			m.filePickerMode = "import"
			m.filePickerPath = filepath.Join(os.Getenv("HOME"))
			m.filePickerPrompt = false
			m.loadFileList()
		case 1: // Export
			m.state = filePickerView
			m.filePickerMode = "export"
			m.filePickerPath = filepath.Join(os.Getenv("HOME"))
			m.filePickerPrompt = false
			ti := textinput.New()
			ti.Placeholder = "ssh-servers-export.json"
			ti.Width = 40
			ti.Prompt = "Filename: "
			m.filePickerInput = ti
			m.loadFileList()
		case 2: // Back
			m.state = listView
		}
		return m, nil

	case "esc":
		m.state = listView
		return m, nil
	}

	return m, nil
}

func (m model) updateFormView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit

	case "esc":
		m.state = listView
		m.message = ""
		return m, nil

	case "p":
		// If on PEM field, open multiline editor
		if m.focusIndex == 5 {
			m.pemBuffer = m.inputs[5].Value()
			prevState := m.state
			m.state = pemEditView
			// Store previous state to return to
			if prevState == addView {
				m.editingID = -1 // -1 means we're adding
			}
			return m, nil
		}

	case "tab", "shift+tab", "up", "down":
		s := msg.String()
		if s == "up" || s == "shift+tab" {
			m.focusIndex--
		} else {
			m.focusIndex++
		}

		if m.focusIndex > len(m.inputs) {
			m.focusIndex = 0
		} else if m.focusIndex < 0 {
			m.focusIndex = len(m.inputs)
		}

		for i := 0; i < len(m.inputs); i++ {
			if i == m.focusIndex {
				m.inputs[i].Focus()
			} else {
				m.inputs[i].Blur()
			}
		}

		return m, nil

	case "enter":
		if m.focusIndex == len(m.inputs) {
			// Submit button pressed
			if m.validateAndSave() {
				m.state = listView
			}
			return m, nil
		}
	}

	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m model) updatePemEditView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit

	case "ctrl+s":
		// Save and return to form
		m.inputs[5].SetValue(m.pemBuffer)
		if m.editingID == -1 {
			m.state = addView
		} else {
			m.state = editView
		}
		m.message = "PEM key saved"
		return m, nil

	case "esc":
		// Cancel and return to form without saving
		if m.editingID == -1 {
			m.state = addView
		} else {
			m.state = editView
		}
		return m, nil

	default:
		// Handle regular text input
		key := msg.String()

		// Handle backspace
		if key == "backspace" {
			if len(m.pemBuffer) > 0 {
				m.pemBuffer = m.pemBuffer[:len(m.pemBuffer)-1]
			}
		} else if key == "enter" {
			m.pemBuffer += "\n"
		} else if len(key) == 1 {
			// Single character input
			m.pemBuffer += key
		}
	}

	return m, nil
}

func (m *model) updateInputs(msg tea.KeyMsg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m *model) validateAndSave() bool {
	name := strings.TrimSpace(m.inputs[0].Value())
	host := strings.TrimSpace(m.inputs[1].Value())
	portStr := strings.TrimSpace(m.inputs[2].Value())
	username := strings.TrimSpace(m.inputs[3].Value())
	password := m.inputs[4].Value()
	pemKey := m.inputs[5].Value()

	if name == "" {
		m.message = "Error: Name is required"
		return false
	}
	if host == "" {
		m.message = "Error: Host is required"
		return false
	}
	if username == "" {
		m.message = "Error: Username is required"
		return false
	}

	// Check if both password and PEM key are provided
	if password != "" && pemKey != "" {
		m.message = "Error: Use either password OR PEM key, not both"
		return false
	}

	port := 22
	if portStr != "" {
		var err error
		port, err = strconv.Atoi(portStr)
		if err != nil || port < 1 || port > 65535 {
			m.message = "Error: Invalid port number"
			return false
		}
	}

	// Validate PEM key format if provided
	if pemKey != "" {
		normalized := normalizePemKey(pemKey)
		if !strings.Contains(normalized, "BEGIN") || !strings.Contains(normalized, "PRIVATE KEY") {
			m.message = "Error: Invalid PEM key format (must include -----BEGIN ... PRIVATE KEY----- and -----END ... PRIVATE KEY-----)"
			return false
		}
	}

	if m.state == addView {
		server := Server{
			ID:       m.config.NextID,
			Name:     name,
			Host:     host,
			Port:     port,
			Username: username,
			Password: password,
			PemKey:   pemKey,
		}
		m.config.Servers = append(m.config.Servers, server)
		m.config.NextID++
		m.message = fmt.Sprintf("Added server: %s", name)
	} else if m.state == editView {
		for i, server := range m.config.Servers {
			if server.ID == m.editingID {
				m.config.Servers[i].Name = name
				m.config.Servers[i].Host = host
				m.config.Servers[i].Port = port
				m.config.Servers[i].Username = username
				m.config.Servers[i].Password = password
				m.config.Servers[i].PemKey = pemKey
				m.message = fmt.Sprintf("Updated server: %s", name)
				break
			}
		}
	}

	if err := m.saveConfig(); err != nil {
		m.message = fmt.Sprintf("Error saving: %v", err)
		return false
	}

	m.refreshList()
	return true
}

func (m *model) deleteServer(id int) {
	newServers := []Server{}
	for _, server := range m.config.Servers {
		if server.ID != id {
			newServers = append(newServers, server)
		}
	}
	m.config.Servers = newServers
	m.saveConfig()
	m.refreshList()
}

func (m *model) refreshList() {
	items := make([]list.Item, len(m.config.Servers))
	for i, server := range m.config.Servers {
		items[i] = server
	}
	m.list.SetItems(items)
}

// normalizePemKey converts escaped newlines, trims surrounding quotes/space,
// normalizes line endings, and ensures a trailing newline.
func normalizePemKey(pem string) string {
	clean := strings.TrimSpace(pem)

	// Handle JSON-style escaped newlines
	clean = strings.ReplaceAll(clean, "\\r\\n", "\n")
	clean = strings.ReplaceAll(clean, "\\n", "\n")

	// Normalize Windows line endings
	clean = strings.ReplaceAll(clean, "\r\n", "\n")

	// Remove surrounding quotes if pasted with them
	clean = strings.Trim(clean, "\"")

	// If the key is on a single line, re-wrap it into a proper PEM block
	if strings.Contains(clean, "-----BEGIN") && strings.Contains(clean, "-----END") && !strings.Contains(clean, "\n") {
		beginIdx := strings.Index(clean, "-----BEGIN")
		endIdx := strings.Index(clean, "-----END")
		if beginIdx >= 0 && endIdx > beginIdx {
			// Find the closing dashes that end the BEGIN header (e.g. '-----BEGIN RSA PRIVATE KEY-----')
			headerCloseRel := strings.Index(clean[beginIdx+len("-----BEGIN"):], "-----")
			if headerCloseRel >= 0 {
				headerEnd := beginIdx + len("-----BEGIN") + headerCloseRel + len("-----")
				header := clean[beginIdx:headerEnd]
				footer := clean[endIdx:]
				body := strings.TrimSpace(clean[headerEnd:endIdx])
				// Remove whitespace characters inside the body only
				body = strings.ReplaceAll(body, " ", "")
				body = strings.ReplaceAll(body, "\t", "")
				body = strings.ReplaceAll(body, "\r", "")
				body = strings.ReplaceAll(body, "\n", "")
				// Wrap body at 64 chars per line
				var wrapped []string
				for i := 0; i < len(body); i += 64 {
					end := i + 64
					if end > len(body) {
						end = len(body)
					}
					wrapped = append(wrapped, body[i:end])
				}
				clean = header + "\n" + strings.Join(wrapped, "\n") + "\n" + footer
			}
		}
	}

	// Ensure trailing newline for OpenSSH parser compatibility
	if clean != "" && !strings.HasSuffix(clean, "\n") {
		clean += "\n"
	}

	return clean
}

func (m *model) connectSSH(server Server) *exec.Cmd {
	args := []string{
		fmt.Sprintf("%s@%s", server.Username, server.Host),
		"-p", strconv.Itoa(server.Port),
	}

	// If PEM key is provided, save it to a temporary file
	if server.PemKey != "" {
		// Create temp directory if it doesn't exist
		tempDir := filepath.Join(os.TempDir(), "termius-from-walmart-keys")
		os.MkdirAll(tempDir, 0700)

		// Create temp key file
		keyFile := filepath.Join(tempDir, fmt.Sprintf("key_%d.pem", server.ID))

		// Clean and normalize the PEM key
		cleanKey := normalizePemKey(server.PemKey)

		// Write the key file
		if err := ioutil.WriteFile(keyFile, []byte(cleanKey), 0600); err == nil {
			// Add SSH options for better compatibility
			args = append([]string{
				"-i", keyFile,
				"-o", "StrictHostKeyChecking=no",
				"-o", "UserKnownHostsFile=/dev/null",
				"-o", "IdentitiesOnly=yes",
			}, args...)
			return exec.Command("ssh", args...)
		}
	}

	// If password is provided, use sshpass
	if server.Password != "" {
		return exec.Command("sshpass", append([]string{"-p", server.Password, "ssh"}, args...)...)
	}

	// Default: use system SSH keys
	return exec.Command("ssh", args...)
}

func (m *model) exportServers() {
	exportPath := filepath.Join(os.Getenv("HOME"), "ssh-servers-export.json")

	exportData := make([]map[string]interface{}, len(m.config.Servers))
	for i, server := range m.config.Servers {
		exportData[i] = map[string]interface{}{
			"name":     server.Name,
			"host":     server.Host,
			"port":     server.Port,
			"username": server.Username,
			"password": server.Password,
			"pem_key":  server.PemKey,
		}
	}

	data, err := json.MarshalIndent(exportData, "", "  ")
	if err != nil {
		m.message = fmt.Sprintf("Export failed: %v", err)
		return
	}

	if err := ioutil.WriteFile(exportPath, data, 0600); err != nil {
		m.message = fmt.Sprintf("Export failed: %v", err)
		return
	}

	m.message = fmt.Sprintf("Exported %d servers to %s", len(m.config.Servers), exportPath)
}

func (m *model) importServers() {
	importPath := filepath.Join(os.Getenv("HOME"), "ssh-servers-import.json")

	data, err := ioutil.ReadFile(importPath)
	if err != nil {
		m.message = fmt.Sprintf("Import failed: %v (place file at %s)", err, importPath)
		return
	}

	var importData []map[string]interface{}
	if err := json.Unmarshal(data, &importData); err != nil {
		m.message = fmt.Sprintf("Import failed: invalid JSON format")
		return
	}

	count := 0
	for _, item := range importData {
		server := Server{
			ID:       m.config.NextID,
			Name:     fmt.Sprintf("%v", item["name"]),
			Host:     fmt.Sprintf("%v", item["host"]),
			Username: fmt.Sprintf("%v", item["username"]),
		}

		if port, ok := item["port"].(float64); ok {
			server.Port = int(port)
		} else {
			server.Port = 22
		}

		if password, ok := item["password"].(string); ok {
			server.Password = password
		}

		if pemKey, ok := item["pem_key"].(string); ok {
			server.PemKey = pemKey
		}

		m.config.Servers = append(m.config.Servers, server)
		m.config.NextID++
		count++
	}

	if err := m.saveConfig(); err != nil {
		m.message = fmt.Sprintf("Import failed: %v", err)
		return
	}

	m.refreshList()
	m.message = fmt.Sprintf("Imported %d servers from %s", count, importPath)
}

// --- File picker helpers ---

type fileItem string

func (f fileItem) FilterValue() string { return string(f) }
func (f fileItem) Title() string       { return string(f) }
func (f fileItem) Description() string { return "" }

func (m *model) loadFileList() {
	entries, err := ioutil.ReadDir(m.filePickerPath)
	if err != nil {
		m.filePickerList.SetItems([]list.Item{})
		m.message = fmt.Sprintf("Unable to read %s: %v", m.filePickerPath, err)
		return
	}

	items := make([]list.Item, 0, len(entries)+1)
	// Parent dir
	if parent := filepath.Dir(m.filePickerPath); parent != m.filePickerPath {
		items = append(items, fileItem("../"))
	}

	for _, e := range entries {
		name := e.Name()
		// skip hidden files/dirs unless toggled
		if !m.filePickerShowHidden && strings.HasPrefix(name, ".") {
			continue
		}
		if e.IsDir() {
			name = name + "/"
		}
		items = append(items, fileItem(name))
	}

	// reuse existing filePickerList to preserve size and other settings
	m.filePickerList.SetItems(items)
	m.filePickerList.Title = fmt.Sprintf("Select file (%s)", m.filePickerMode)
}

func (m model) updateFilePickerView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// If typing filename for export
	if m.filePickerPrompt {
		// allow textinput to handle the key
		var cmd tea.Cmd
		m.filePickerInput, cmd = m.filePickerInput.Update(msg)
		if msg.String() == "enter" {
			filename := strings.TrimSpace(m.filePickerInput.Value())
			if filename != "" {
				full := filepath.Join(m.filePickerPath, filename)
				m.exportServersToPath(full)
				m.filePickerPrompt = false
				m.state = listView
			} else {
				m.message = "Filename cannot be empty"
			}
		} else if msg.String() == "esc" {
			m.filePickerPrompt = false
		}
		return m, cmd
	}

	switch msg.String() {
	case "ctrl+c", "q", "esc":
		m.state = listView
		return m, nil
	case ".":
		// toggle hidden files
		m.filePickerShowHidden = !m.filePickerShowHidden
		m.loadFileList()
		return m, nil
	case "enter":
		sel := m.filePickerList.SelectedItem()
		if sel == nil {
			return m, nil
		}
		name := sel.FilterValue()
		// handle parent
		if name == "../" {
			parent := filepath.Dir(m.filePickerPath)
			m.filePickerPath = parent
			m.loadFileList()
			return m, nil
		}

		// directory?
		if strings.HasSuffix(name, "/") {
			// descend
			dirName := strings.TrimSuffix(name, "/")
			m.filePickerPath = filepath.Join(m.filePickerPath, dirName)
			m.loadFileList()
			return m, nil
		}

		// file chosen
		full := filepath.Join(m.filePickerPath, name)
		if m.filePickerMode == "import" {
			m.importServersFromPath(full)
			m.state = listView
			return m, nil
		}
		// export mode: export to selected file (overwrite)
		if m.filePickerMode == "export" {
			m.exportServersToPath(full)
			m.state = listView
			return m, nil
		}

	case "x":
		// quick export to current dir using default name
		if m.filePickerMode == "export" {
			full := filepath.Join(m.filePickerPath, "ssh-servers-export.json")
			m.exportServersToPath(full)
			m.state = listView
			return m, nil
		}
	case "n":
		// new filename for export
		if m.filePickerMode == "export" {
			m.filePickerPrompt = true
			m.filePickerInput.SetValue("")
			m.filePickerInput.Focus()
			return m, nil
		}
	}

	// delegate list navigation
	var cmd tea.Cmd
	m.filePickerList, cmd = m.filePickerList.Update(msg)
	return m, cmd
}

// compact delegate for file items (single-line, no extra spacing)
type fileDelegate struct{}

func (d fileDelegate) Height() int                             { return 1 }
func (d fileDelegate) Spacing() int                            { return 0 }
func (d fileDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d fileDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	fi, ok := listItem.(fileItem)
	if !ok {
		return
	}
	name := string(fi)

	fn := fileItemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return fileSelectedStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(name))
}

func (m model) viewFilePicker() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render(fmt.Sprintf("File Picker — %s", m.filePickerMode)) + "\n\n")
	b.WriteString(helpStyle.Render("Path: ") + m.filePickerPath + "\n\n")

	if m.filePickerPrompt {
		b.WriteString(m.filePickerInput.View() + "\n\n")
		b.WriteString(helpStyle.Render("Type filename and press Enter to save, Esc to cancel"))
		return b.String()
	}

	b.WriteString(m.filePickerList.View() + "\n\n")
	if m.filePickerMode == "import" {
		b.WriteString(helpStyle.Render("Enter: open file / enter dir • Esc: cancel • .: toggle hidden"))
	} else {
		b.WriteString(helpStyle.Render("Enter: choose file (overwrite) • x: export here • n: new filename • Esc: cancel • .: toggle hidden"))
	}

	if m.message != "" {
		msgStyle := messageStyle
		if strings.HasPrefix(m.message, "Error") {
			msgStyle = errorStyle
		}
		b.WriteString("\n\n" + msgStyle.Render(m.message))
	}

	return b.String()
}

func (m *model) importServersFromPath(importPath string) {
	data, err := ioutil.ReadFile(importPath)
	if err != nil {
		m.message = fmt.Sprintf("Import failed: %v (path %s)", err, importPath)
		return
	}

	var importData []map[string]interface{}
	if err := json.Unmarshal(data, &importData); err != nil {
		m.message = "Import failed: invalid JSON format"
		return
	}

	count := 0
	for _, item := range importData {
		server := Server{
			ID:       m.config.NextID,
			Name:     fmt.Sprintf("%v", item["name"]),
			Host:     fmt.Sprintf("%v", item["host"]),
			Username: fmt.Sprintf("%v", item["username"]),
		}

		if port, ok := item["port"].(float64); ok {
			server.Port = int(port)
		} else {
			server.Port = 22
		}

		if password, ok := item["password"].(string); ok {
			server.Password = password
		}

		if pemKey, ok := item["pem_key"].(string); ok {
			server.PemKey = pemKey
		}

		m.config.Servers = append(m.config.Servers, server)
		m.config.NextID++
		count++
	}

	if err := m.saveConfig(); err != nil {
		m.message = fmt.Sprintf("Import failed: %v", err)
		return
	}

	m.refreshList()
	m.message = fmt.Sprintf("Imported %d servers from %s", count, importPath)
}

func (m *model) exportServersToPath(exportPath string) {
	exportData := make([]map[string]interface{}, len(m.config.Servers))
	for i, server := range m.config.Servers {
		exportData[i] = map[string]interface{}{
			"name":     server.Name,
			"host":     server.Host,
			"port":     server.Port,
			"username": server.Username,
			"password": server.Password,
			"pem_key":  server.PemKey,
		}
	}

	data, err := json.MarshalIndent(exportData, "", "  ")
	if err != nil {
		m.message = fmt.Sprintf("Export failed: %v", err)
		return
	}

	if err := ioutil.WriteFile(exportPath, data, 0600); err != nil {
		m.message = fmt.Sprintf("Export failed: %v", err)
		return
	}

	m.message = fmt.Sprintf("Exported %d servers to %s", len(m.config.Servers), exportPath)
}

func (m model) View() string {
	switch m.state {
	case listView:
		return m.viewList()
	case addView:
		return m.viewForm("Add New Server")
	case editView:
		return m.viewForm("Edit Server")
	case menuView:
		return m.viewMenu()
	case filePickerView:
		return m.viewFilePicker()
	case pemEditView:
		return m.viewPemEdit()
	}
	return ""
}

func (m model) viewList() string {
	help := helpStyle.Render("\nKeys: [a]dd • [e]dit • [d]elete • [enter] connect • [m]enu • [q]uit")

	if m.message != "" {
		msgStyle := messageStyle
		if strings.HasPrefix(m.message, "Error") {
			msgStyle = errorStyle
		}
		return m.list.View() + "\n" + msgStyle.Render(m.message) + help
	}

	return m.list.View() + help
}

func (m model) viewMenu() string {
	s := titleStyle.Render("Import/Export Menu") + "\n\n"

	for i, option := range m.menuOptions {
		cursor := " "
		if m.menuCursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, option)
	}

	s += "\n" + helpStyle.Render("Use ↑/↓ to navigate, [enter] to select, [esc] to go back")

	if m.message != "" {
		msgStyle := messageStyle
		if strings.HasPrefix(m.message, "Error") || strings.HasPrefix(m.message, "Import failed") || strings.HasPrefix(m.message, "Export failed") {
			msgStyle = errorStyle
		}
		s += "\n\n" + msgStyle.Render(m.message)
	}

	return s
}

func (m model) viewForm(title string) string {
	var b strings.Builder

	b.WriteString(titleStyle.Render(title) + "\n\n")

	for i, input := range m.inputs {
		b.WriteString(input.View())

		// Add hint for PEM field
		if i == 5 && m.focusIndex == 5 {
			b.WriteString(helpStyle.Render(" (Press 'p' for multiline editor)"))
		}

		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := "[Submit]"
	if m.focusIndex == len(m.inputs) {
		button = "> [Submit] <"
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", button)

	b.WriteString(helpStyle.Render("Navigate: [tab]/[shift+tab] • Submit: [enter] • Cancel: [esc]"))

	if m.message != "" {
		msgStyle := messageStyle
		if strings.HasPrefix(m.message, "Error") {
			msgStyle = errorStyle
		}
		b.WriteString("\n\n" + msgStyle.Render(m.message))
	}

	return b.String()
}

func (m model) viewPemEdit() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("PEM Key Editor") + "\n\n")
	b.WriteString(helpStyle.Render("Paste your PEM private key below:") + "\n\n")

	// Show the PEM buffer with a border
	pemDisplay := m.pemBuffer
	if pemDisplay == "" {
		pemDisplay = "(empty - paste your PEM key here)"
	}

	// Simple box around the PEM content
	lines := strings.Split(pemDisplay, "\n")
	maxLen := 70

	b.WriteString("┌" + strings.Repeat("─", maxLen) + "┐\n")
	for _, line := range lines {
		if len(line) > maxLen-2 {
			line = line[:maxLen-2]
		}
		padding := maxLen - len(line)
		b.WriteString("│ " + line + strings.Repeat(" ", padding-1) + "│\n")
	}
	// Add some empty lines for visual space
	for i := len(lines); i < 15; i++ {
		b.WriteString("│" + strings.Repeat(" ", maxLen) + "│\n")
	}
	b.WriteString("└" + strings.Repeat("─", maxLen) + "┘\n\n")

	b.WriteString(helpStyle.Render("Lines: "+strconv.Itoa(len(lines))) + "\n")
	b.WriteString(helpStyle.Render("Characters: "+strconv.Itoa(len(m.pemBuffer))) + "\n\n")

	b.WriteString(messageStyle.Render("[Ctrl+S] Save • [Esc] Cancel") + "\n")

	if m.message != "" {
		msgStyle := messageStyle
		if strings.HasPrefix(m.message, "Error") {
			msgStyle = errorStyle
		}
		b.WriteString("\n" + msgStyle.Render(m.message))
	}

	return b.String()
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
