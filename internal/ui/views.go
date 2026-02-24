package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ── App Header ─────────────────────────────────────────────────────────

const appName = "Termius from Walmart"

func (m Model) renderHeader() string {
	logo := AppHeaderStyle.Render("  " + appName)
	count := fmt.Sprintf(" %d servers", len(m.Config.Servers))
	info := AppVersionStyle.Render(count)
	return logo + " " + info + "\n"
}

// ── Key Bar Helpers ────────────────────────────────────────────────────

func renderKeyBar(width int, pairs ...string) string {
	var parts []string
	for i := 0; i < len(pairs)-1; i += 2 {
		parts = append(parts, RenderKey(pairs[i], pairs[i+1]))
	}
	bar := strings.Join(parts, "    ")
	return StatusBarStyle.Width(width).Render(bar)
}

// ── List View ──────────────────────────────────────────────────────────

// viewList renders the main server list screen.
func (m Model) viewList() string {
	var b strings.Builder

	b.WriteString(m.renderHeader())
	b.WriteString(m.List.View())

	// Message area
	if m.Message != "" {
		msgStyle := MessageStyle
		if strings.HasPrefix(m.Message, "Error") {
			msgStyle = ErrorStyle
		}
		b.WriteString("\n " + msgStyle.Render(m.Message))
	}

	// Key bar
	b.WriteString("\n")
	w := m.Width
	if w == 0 {
		w = 80
	}
	b.WriteString(renderKeyBar(w,
		"a", "add", "e", "edit", "d", "delete",
		"enter", "connect", "s", "sftp", "m", "menu", "?", "help", "q", "quit",
	))

	return b.String()
}

// ── Menu View ──────────────────────────────────────────────────────────

// viewMenu renders the import/export menu.
func (m Model) viewMenu() string {
	var b strings.Builder

	b.WriteString(m.renderHeader())
	b.WriteString(TitleStyle.Render("Import / Export") + "\n\n")

	menuIcons := []string{"  ", "  ", "  "}
	for i, option := range m.MenuOptions {
		if m.MenuCursor == i {
			icon := MenuIconStyle.Render(menuIcons[i])
			b.WriteString(MenuItemSelectedStyle.Render("> "+icon+option) + "\n")
		} else {
			icon := lipgloss.NewStyle().Foreground(ColorTextMuted).Render(menuIcons[i])
			b.WriteString(MenuItemStyle.Render("  "+icon+option) + "\n")
		}
	}

	if m.Message != "" {
		msgStyle := MessageStyle
		if strings.HasPrefix(m.Message, "Error") || strings.HasPrefix(m.Message, "Import failed") || strings.HasPrefix(m.Message, "Export failed") {
			msgStyle = ErrorStyle
		}
		b.WriteString("\n" + msgStyle.Render(" "+m.Message))
	}

	b.WriteString("\n\n")
	w := m.Width
	if w == 0 {
		w = 80
	}
	b.WriteString(renderKeyBar(w, "enter", "select", "esc", "back"))

	return b.String()
}

// ── Form View ──────────────────────────────────────────────────────────

// viewForm renders the add/edit server form.
func (m Model) viewForm(title string) string {
	var b strings.Builder

	b.WriteString(m.renderHeader())
	b.WriteString(TitleStyle.Render(title) + "\n\n")

	labels := []string{"Name", "Host", "Port", "User", "Pass", "PEM", "SFTP"}
	for i, input := range m.Inputs {
		label := InputLabelStyle.Render(labels[i])
		var field string
		if i == m.FocusIndex {
			field = InputActiveStyle.Render(input.View())
		} else {
			field = InputInactiveStyle.Render(input.View())
		}
		row := lipgloss.JoinHorizontal(lipgloss.Center, label, field)
		b.WriteString(row)
		if i == 5 && m.FocusIndex == 5 {
			b.WriteString(HelpStyle.Render("  press p for multiline editor"))
		}
		b.WriteString("\n")
	}

	// Submit button
	if m.FocusIndex == len(m.Inputs) {
		b.WriteString("\n" + lipgloss.NewStyle().MarginLeft(9).Render(
			SubmitButtonActiveStyle.Render("  Submit  "),
		))
	} else {
		b.WriteString("\n" + lipgloss.NewStyle().MarginLeft(9).Render(
			SubmitButtonInactiveStyle.Render("  Submit  "),
		))
	}

	// Message
	if m.Message != "" {
		msgStyle := MessageStyle
		if strings.HasPrefix(m.Message, "Error") {
			msgStyle = ErrorStyle
		}
		b.WriteString("\n\n " + msgStyle.Render(m.Message))
	}

	b.WriteString("\n\n")
	w := m.Width
	if w == 0 {
		w = 80
	}
	b.WriteString(renderKeyBar(w, "tab", "next", "shift+tab", "prev", "enter", "submit", "esc", "cancel"))

	return b.String()
}

// ── PEM Editor View ────────────────────────────────────────────────────

// viewPemEdit renders the PEM key multiline editor.
func (m Model) viewPemEdit() string {
	var b strings.Builder

	b.WriteString(m.renderHeader())
	b.WriteString(TitleStyle.Render("PEM Key Editor") + "\n")
	b.WriteString(SubtitleStyle.Render("Paste your PEM private key below") + "\n\n")

	pemDisplay := m.PemBuffer
	if pemDisplay == "" {
		pemDisplay = "(empty - paste your PEM key here)"
	}

	lines := strings.Split(pemDisplay, "\n")
	maxLen := 68

	// Build the PEM content with line numbers
	var pemContent strings.Builder
	for i, line := range lines {
		if len(line) > maxLen-6 {
			line = line[:maxLen-6]
		}
		lineNum := PemLineNumStyle.Render(fmt.Sprintf("%d", i+1))
		content := PemContentStyle.Render(line)
		pemContent.WriteString(lineNum + " " + content + "\n")
	}
	// Pad empty lines to minimum height
	for i := len(lines); i < 12; i++ {
		lineNum := PemLineNumStyle.Render(fmt.Sprintf("%d", i+1))
		pemContent.WriteString(lineNum + "\n")
	}

	b.WriteString(PemBoxStyle.Render(pemContent.String()))

	// Stats
	b.WriteString("\n")
	stats := HelpStyle.Render(fmt.Sprintf("  Lines: %d  Characters: %d", len(lines), len(m.PemBuffer)))
	b.WriteString(stats)

	b.WriteString("\n\n")
	w := m.Width
	if w == 0 {
		w = 80
	}
	b.WriteString(renderKeyBar(w, "ctrl+s", "save", "esc", "cancel"))

	return b.String()
}

// ── File Picker View ───────────────────────────────────────────────────

// viewFilePicker renders the file picker view.
func (m Model) viewFilePicker() string {
	var b strings.Builder

	b.WriteString(m.renderHeader())
	title := "File Picker"
	if m.FilePickerMode == "import" {
		title = "Import File"
	} else {
		title = "Export Location"
	}
	b.WriteString(TitleStyle.Render(title) + "\n")
	b.WriteString(InfoStyle.Render("  "+m.FilePickerPath) + "\n\n")

	if m.FilePickerPrompt {
		label := InputLabelStyle.Render("File")
		field := InputActiveStyle.Render(m.FilePickerInput.View())
		b.WriteString(lipgloss.JoinHorizontal(lipgloss.Center, label, field) + "\n\n")
		w := m.Width
		if w == 0 {
			w = 80
		}
		b.WriteString(renderKeyBar(w, "enter", "save", "esc", "cancel"))
		return b.String()
	}

	b.WriteString(m.FilePickerList.View() + "\n")

	if m.Message != "" {
		msgStyle := MessageStyle
		if strings.HasPrefix(m.Message, "Error") {
			msgStyle = ErrorStyle
		}
		b.WriteString("\n " + msgStyle.Render(m.Message))
	}

	b.WriteString("\n")
	w := m.Width
	if w == 0 {
		w = 80
	}
	if m.FilePickerMode == "import" {
		b.WriteString(renderKeyBar(w, "enter", "open", ".", "toggle hidden", "esc", "cancel"))
	} else {
		b.WriteString(renderKeyBar(w, "enter", "overwrite", "x", "export here", "n", "new name", ".", "toggle hidden", "esc", "cancel"))
	}

	return b.String()
}

// ── SFTP View ──────────────────────────────────────────────────────────

// viewSFTP renders the SFTP split-screen view.
func (m Model) viewSFTP() string {
	server := m.SelectedServer
	if server == nil {
		return ""
	}

	var b strings.Builder

	// Header
	b.WriteString(m.renderHeader())
	connInfo := InfoStyle.Render(fmt.Sprintf("  Connected: %s@%s:%d", server.Username, server.Host, server.Port))
	b.WriteString(connInfo + "\n")

	// Transfer progress
	if m.IsTransferring {
		bar := TransferBarStyle.Render(fmt.Sprintf("  [%d%%] %s", m.TransferProgress, m.TransferMessage))
		b.WriteString(bar + "\n")
	}
	b.WriteString("\n")

	// Calculate pane width
	totalWidth := m.Width
	if totalWidth == 0 {
		totalWidth = 80
	}
	paneWidth := (totalWidth - 5) / 2 // 5 = divider + margins
	if paneWidth < MinPaneWidth {
		paneWidth = MinPaneWidth
	}

	// Local pane
	var localHeader, remoteHeader string
	var localPanel, remotePanel lipgloss.Style

	if m.FocusPane == "local" {
		localHeader = SFTPHeaderStyle.Width(paneWidth - 4).Render("LOCAL")
		remoteHeader = SFTPHeaderDimStyle.Width(paneWidth - 4).Render("REMOTE")
		localPanel = PanelFocusedStyle.Width(paneWidth)
		remotePanel = PanelDimStyle.Width(paneWidth)
	} else {
		localHeader = SFTPHeaderDimStyle.Width(paneWidth - 4).Render("LOCAL")
		remoteHeader = SFTPHeaderStyle.Width(paneWidth - 4).Render("REMOTE")
		localPanel = PanelDimStyle.Width(paneWidth)
		remotePanel = PanelFocusedStyle.Width(paneWidth)
	}

	localPath := SFTPPathStyle.Width(paneWidth - 4).Render(truncatePath(m.LocalPath, paneWidth-6))
	remotePath := SFTPPathStyle.Width(paneWidth - 4).Render(truncatePath(m.RemotePath, paneWidth-6))

	localContent := localHeader + "\n" + localPath + "\n\n" + m.LocalFileList.View()
	remoteContent := remoteHeader + "\n" + remotePath + "\n\n" + m.RemoteFileList.View()

	leftPane := localPanel.Render(localContent)
	rightPane := remotePanel.Render(remoteContent)

	divider := SFTPDividerStyle.Render(" ")

	b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, leftPane, divider, rightPane))

	// Message
	if m.Message != "" {
		msgStyle := MessageStyle
		if strings.HasPrefix(m.Message, "Error") {
			msgStyle = ErrorStyle
		}
		b.WriteString("\n " + msgStyle.Render(m.Message))
	}

	b.WriteString("\n")
	b.WriteString(renderKeyBar(totalWidth,
		"tab", "switch pane", "c", "copy", "d", "delete", "enter", "navigate", "?", "help", "q", "quit",
	))

	return b.String()
}

// ── Confirmation Dialog ────────────────────────────────────────────────

// viewConfirm renders the confirmation dialog as a centered modal.
func (m Model) viewConfirm() string {
	var content strings.Builder

	content.WriteString(ConfirmTitleStyle.Render("  "+m.ConfirmTitle) + "\n\n")
	content.WriteString(ConfirmMessageStyle.Render(m.ConfirmMessage) + "\n\n")

	var noBtn, yesBtn string
	if m.ConfirmCursor == 0 {
		noBtn = ButtonActiveStyle.Render("  No  ")
		yesBtn = ButtonInactiveStyle.Render("  Yes  ")
	} else {
		noBtn = ButtonInactiveStyle.Render("  No  ")
		yesBtn = ButtonDangerActiveStyle.Render("  Yes  ")
	}

	buttons := lipgloss.JoinHorizontal(lipgloss.Center, noBtn, "   ", yesBtn)
	content.WriteString(lipgloss.NewStyle().Align(lipgloss.Center).Width(46).Render(buttons))
	content.WriteString("\n\n")
	content.WriteString(HelpStyle.Render("arrow keys to select  enter to confirm  esc to cancel"))

	modal := ModalStyle.Render(content.String())

	// Center the modal on screen
	if m.Width > 0 && m.Height > 0 {
		return lipgloss.Place(m.Width, m.Height,
			lipgloss.Center, lipgloss.Center,
			modal)
	}
	return "\n\n" + modal
}

// ── Confirm View Update ────────────────────────────────────────────────

// updateConfirmView handles key events in the confirmation dialog.
func (m Model) updateConfirmView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		m.State = m.PreviousState
		if m.State == ConfirmView || m.State == 0 {
			m.State = ListView
		}
		return m, nil
	case "left", "h":
		m.ConfirmCursor = 0
	case "right", "l":
		m.ConfirmCursor = 1
	case "enter":
		if m.ConfirmCursor == 1 && m.ConfirmAction != nil {
			m.ConfirmAction()
		}
		m.State = m.PreviousState
		if m.State == ConfirmView || m.State == 0 {
			m.State = ListView
		}
		return m, nil
	}
	return m, nil
}

// ── Helpers ────────────────────────────────────────────────────────────

func truncatePath(path string, maxLen int) string {
	if len(path) <= maxLen {
		return path
	}
	return "..." + path[len(path)-maxLen+3:]
}

// ── Connecting View ────────────────────────────────────────────────────

var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// viewConnecting renders the spinner while establishing SFTP connection.
func (m Model) viewConnecting() string {
	var b strings.Builder

	b.WriteString(m.renderHeader())

	frame := spinnerFrames[m.SpinnerFrame%len(spinnerFrames)]
	spinner := lipgloss.NewStyle().Foreground(ColorAccent).Bold(true).Render(frame)

	serverName := m.ConnectingName
	if m.SelectedServer != nil {
		serverName = fmt.Sprintf("%s (%s@%s:%d)",
			m.SelectedServer.Name, m.SelectedServer.Username,
			m.SelectedServer.Host, m.SelectedServer.Port)
	}

	msg := lipgloss.NewStyle().Foreground(ColorText).Render(
		fmt.Sprintf(" Connecting to %s...", serverName))

	content := "\n\n" + spinner + msg + "\n\n" +
		HelpStyle.Render("  Establishing SFTP connection, please wait...")

	if m.Width > 0 && m.Height > 0 {
		box := PanelStyle.Width(ClampWidth(m.Width-10, MinPaneWidth, 70)).Render(content)
		return b.String() + "\n" + lipgloss.Place(m.Width, m.Height-4,
			lipgloss.Center, lipgloss.Center, box)
	}
	b.WriteString(PanelStyle.Render(content))
	return b.String()
}

// updateConnectingView handles keys during the connecting spinner.
func (m Model) updateConnectingView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		// Cancel connecting — go back to list
		m.State = ListView
		m.SelectedServer = nil
		m.ConnectingName = ""
		m.Message = "Connection cancelled"
		m.MessageTimer++
		return m, clearMessageAfter(3 * secondDuration)
	}
	return m, nil
}

// ── Help Overlay ───────────────────────────────────────────────────────

// viewHelp renders a keyboard shortcut help overlay.
func (m Model) viewHelp() string {
	var b strings.Builder

	b.WriteString(m.renderHeader())

	var content strings.Builder
	content.WriteString(ConfirmTitleStyle.Render("  Keyboard Shortcuts") + "\n\n")

	// Organized by view
	sections := []struct {
		title string
		keys  []struct{ key, desc string }
	}{
		{
			title: "Server List",
			keys: []struct{ key, desc string }{
				{"a", "Add new server"},
				{"e", "Edit selected server"},
				{"d", "Delete selected server"},
				{"enter", "SSH connect to server"},
				{"s", "Open SFTP file browser"},
				{"m", "Import/Export menu"},
				{"/", "Filter servers"},
				{"q", "Quit application"},
			},
		},
		{
			title: "SFTP Browser",
			keys: []struct{ key, desc string }{
				{"tab", "Switch local/remote pane"},
				{"enter", "Navigate into directory"},
				{"c", "Copy file to other pane"},
				{"d", "Delete selected file"},
				{"q/esc", "Close SFTP & return"},
			},
		},
		{
			title: "Forms & Editors",
			keys: []struct{ key, desc string }{
				{"tab", "Next field"},
				{"shift+tab", "Previous field"},
				{"p", "Open PEM editor (on PEM field)"},
				{"ctrl+s", "Save PEM key"},
				{"enter", "Submit form"},
				{"esc", "Cancel / go back"},
			},
		},
		{
			title: "Global",
			keys: []struct{ key, desc string }{
				{"?", "Toggle this help overlay"},
				{"ctrl+c", "Force quit"},
			},
		},
	}

	for _, section := range sections {
		content.WriteString(SectionHeaderStyle.Render(section.title) + "\n")
		for _, k := range section.keys {
			key := KeyStyle.Render(fmt.Sprintf("%-12s", k.key))
			desc := lipgloss.NewStyle().Foreground(ColorTextDim).Render(k.desc)
			content.WriteString("  " + key + " " + desc + "\n")
		}
		content.WriteString("\n")
	}

	content.WriteString(HelpStyle.Render("Press ? or esc to close"))

	modal := ModalStyle.Width(55).Render(content.String())

	if m.Width > 0 && m.Height > 0 {
		return b.String() + lipgloss.Place(m.Width, m.Height-2,
			lipgloss.Center, lipgloss.Center, modal)
	}
	return b.String() + "\n\n" + modal
}

// updateHelpView handles keys in the help overlay.
func (m Model) updateHelpView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "?", "esc", "q":
		m.State = m.PreviousState
		if m.State == HelpView || m.State == 0 {
			m.State = ListView
		}
		return m, nil
	}
	return m, nil
}
