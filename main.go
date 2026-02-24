package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"termius-from-walmart/internal/ui"
)

// Version is set via ldflags at build time
var Version = "dev"

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Printf("termius-from-walmart %s\n", Version)
		os.Exit(0)
	}

	p := tea.NewProgram(ui.NewModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
