package ui

import "github.com/charmbracelet/lipgloss"

var (
	TitleStyle = lipgloss.NewStyle().
			MarginLeft(2).
			Foreground(lipgloss.Color("62")).
			Bold(true)

	HelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginLeft(2)

	MessageStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("170")).
			MarginLeft(2).
			Bold(true)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			MarginLeft(2).
			Bold(true)

	FileItemStyle     = lipgloss.NewStyle().PaddingLeft(2)
	FileSelectedStyle = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("170"))
)
