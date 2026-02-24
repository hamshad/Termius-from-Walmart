package ui

import (
	"fmt"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// viewList renders the main server list screen.
func (m Model) viewList() string {
	help := HelpStyle.Render("\nKeys: [a]dd  [e]dit  [d]elete  [enter] connect  [s]ftp  [m]enu  [q]uit")

	if m.Message != "" {
		msgStyle := MessageStyle
		if strings.HasPrefix(m.Message, "Error") {
			msgStyle = ErrorStyle
		}
		return m.List.View() + "\n" + msgStyle.Render(m.Message) + help
	}

	return m.List.View() + help
}

// viewMenu renders the import/export menu.
func (m Model) viewMenu() string {
	s := TitleStyle.Render("Import/Export Menu") + "\n\n"

	for i, option := range m.MenuOptions {
		cursor := " "
		if m.MenuCursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, option)
	}

	s += "\n" + HelpStyle.Render("Use arrow keys to navigate, [enter] to select, [esc] to go back")

	if m.Message != "" {
		msgStyle := MessageStyle
		if strings.HasPrefix(m.Message, "Error") || strings.HasPrefix(m.Message, "Import failed") || strings.HasPrefix(m.Message, "Export failed") {
			msgStyle = ErrorStyle
		}
		s += "\n\n" + msgStyle.Render(m.Message)
	}

	return s
}

// viewForm renders the add/edit server form.
func (m Model) viewForm(title string) string {
	var b strings.Builder

	b.WriteString(TitleStyle.Render(title) + "\n\n")

	for i, input := range m.Inputs {
		b.WriteString(input.View())
		if i == 5 && m.FocusIndex == 5 {
			b.WriteString(HelpStyle.Render(" (Press 'p' for multiline editor)"))
		}
		if i < len(m.Inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := "[Submit]"
	if m.FocusIndex == len(m.Inputs) {
		button = "> [Submit] <"
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", button)

	b.WriteString(HelpStyle.Render("Navigate: [tab]/[shift+tab]  Submit: [enter]  Cancel: [esc]"))

	if m.Message != "" {
		msgStyle := MessageStyle
		if strings.HasPrefix(m.Message, "Error") {
			msgStyle = ErrorStyle
		}
		b.WriteString("\n\n" + msgStyle.Render(m.Message))
	}

	return b.String()
}

// viewPemEdit renders the PEM key multiline editor.
func (m Model) viewPemEdit() string {
	var b strings.Builder

	b.WriteString(TitleStyle.Render("PEM Key Editor") + "\n\n")
	b.WriteString(HelpStyle.Render("Paste your PEM private key below:") + "\n\n")

	pemDisplay := m.PemBuffer
	if pemDisplay == "" {
		pemDisplay = "(empty - paste your PEM key here)"
	}

	lines := strings.Split(pemDisplay, "\n")
	maxLen := 70

	b.WriteString("+" + strings.Repeat("-", maxLen) + "+\n")
	for _, line := range lines {
		if len(line) > maxLen-2 {
			line = line[:maxLen-2]
		}
		padding := maxLen - len(line)
		b.WriteString("| " + line + strings.Repeat(" ", padding-1) + "|\n")
	}
	for i := len(lines); i < 15; i++ {
		b.WriteString("|" + strings.Repeat(" ", maxLen) + "|\n")
	}
	b.WriteString("+" + strings.Repeat("-", maxLen) + "+\n\n")

	b.WriteString(HelpStyle.Render("Lines: "+strconv.Itoa(len(lines))) + "\n")
	b.WriteString(HelpStyle.Render("Characters: "+strconv.Itoa(len(m.PemBuffer))) + "\n\n")

	b.WriteString(MessageStyle.Render("[Ctrl+S] Save  [Esc] Cancel") + "\n")

	if m.Message != "" {
		msgStyle := MessageStyle
		if strings.HasPrefix(m.Message, "Error") {
			msgStyle = ErrorStyle
		}
		b.WriteString("\n" + msgStyle.Render(m.Message))
	}

	return b.String()
}

// viewFilePicker renders the file picker view.
func (m Model) viewFilePicker() string {
	var b strings.Builder
	b.WriteString(TitleStyle.Render(fmt.Sprintf("File Picker - %s", m.FilePickerMode)) + "\n\n")
	b.WriteString(HelpStyle.Render("Path: ") + m.FilePickerPath + "\n\n")

	if m.FilePickerPrompt {
		b.WriteString(m.FilePickerInput.View() + "\n\n")
		b.WriteString(HelpStyle.Render("Type filename and press Enter to save, Esc to cancel"))
		return b.String()
	}

	b.WriteString(m.FilePickerList.View() + "\n\n")
	if m.FilePickerMode == "import" {
		b.WriteString(HelpStyle.Render("Enter: open file / enter dir  Esc: cancel  .: toggle hidden"))
	} else {
		b.WriteString(HelpStyle.Render("Enter: choose file (overwrite)  x: export here  n: new filename  Esc: cancel  .: toggle hidden"))
	}

	if m.Message != "" {
		msgStyle := MessageStyle
		if strings.HasPrefix(m.Message, "Error") {
			msgStyle = ErrorStyle
		}
		b.WriteString("\n\n" + msgStyle.Render(m.Message))
	}

	return b.String()
}

// viewSFTP renders the SFTP split-screen view.
func (m Model) viewSFTP() string {
	var b strings.Builder
	server := m.SelectedServer
	if server == nil {
		return ""
	}

	b.WriteString(TitleStyle.Render(fmt.Sprintf("SFTP: %s@%s", server.Username, server.Host)) + "\n")

	if m.IsTransferring {
		progressBar := fmt.Sprintf("[%d%%]", m.TransferProgress)
		b.WriteString(MessageStyle.Render(progressBar) + " " + m.TransferMessage + "\n")
	}
	b.WriteString("\n")

	localHeader := "LOCAL"
	remoteHeader := "REMOTE"
	if m.FocusPane == "local" {
		localHeader = "> " + localHeader + " <"
	} else {
		remoteHeader = "> " + remoteHeader + " <"
	}

	b.WriteString(HelpStyle.Render(localHeader) + "  " + HelpStyle.Render(remoteHeader) + "\n")
	b.WriteString(HelpStyle.Render(m.LocalPath) + " | " + HelpStyle.Render(m.RemotePath) + "\n\n")

	localView := m.LocalFileList.View()
	remoteView := m.RemoteFileList.View()

	localLines := strings.Split(localView, "\n")
	remoteLines := strings.Split(remoteView, "\n")

	maxLines := len(localLines)
	if len(remoteLines) > maxLines {
		maxLines = len(remoteLines)
	}

	for i := 0; i < maxLines; i++ {
		var localLine, remoteLine string
		if i < len(localLines) {
			localLine = localLines[i]
		}
		if i < len(remoteLines) {
			remoteLine = remoteLines[i]
		}
		if len(localLine) < 40 {
			localLine += strings.Repeat(" ", 40-len(localLine))
		}
		b.WriteString(localLine + " | " + remoteLine + "\n")
	}

	b.WriteString("\n" + HelpStyle.Render("Keys: [Tab] switch pane  [c]opy  [d]elete  [enter] navigate  [q]uit"))

	if m.Message != "" {
		msgStyle := MessageStyle
		if strings.HasPrefix(m.Message, "Error") {
			msgStyle = ErrorStyle
		}
		b.WriteString("\n" + msgStyle.Render(m.Message))
	}

	return b.String()
}

// viewConfirm renders the confirmation dialog.
func (m Model) viewConfirm() string {
	var b strings.Builder

	b.WriteString("\n\n")
	b.WriteString(TitleStyle.Render(m.ConfirmTitle) + "\n\n")
	b.WriteString(HelpStyle.Render(m.ConfirmMessage) + "\n\n")

	noBtn := "  [ No ]  "
	yesBtn := "  [ Yes ]  "
	if m.ConfirmCursor == 0 {
		noBtn = "  > [ No ] <  "
	} else {
		yesBtn = "  > [ Yes ] <  "
	}

	b.WriteString("    " + noBtn + "    " + yesBtn + "\n\n")
	b.WriteString(HelpStyle.Render("Use arrow keys to select, [enter] to confirm, [esc] to cancel"))

	return b.String()
}

// updateConfirmView handles key events in the confirmation dialog.
func (m Model) updateConfirmView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		m.State = ListView
		return m, nil
	case "left", "h":
		m.ConfirmCursor = 0
	case "right", "l":
		m.ConfirmCursor = 1
	case "enter":
		if m.ConfirmCursor == 1 && m.ConfirmAction != nil {
			m.ConfirmAction()
		}
		m.State = ListView
		return m, nil
	}
	return m, nil
}
