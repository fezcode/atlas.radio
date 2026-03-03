package main

import (
	"fmt"
	"os"

	"atlas.radio/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

var Version = "dev"

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "-v" || os.Args[1] == "--version") {
		fmt.Printf("atlas.radio v%s\n", Version)
		return
	}
	if len(os.Args) > 1 && (os.Args[1] == "-h" || os.Args[1] == "--help" || os.Args[1] == "help") {
		fmt.Println("Atlas Radio - Global terminal radio receiver with Pip-Boy TUI.")
		fmt.Println("\nUsage:")
		fmt.Println("  atlas.radio        Start the radio receiver")
		fmt.Println("  atlas.radio -v     Show version")
		fmt.Println("  atlas.radio -h     Show this help")
		return
	}

	m := ui.NewModel()
	// Ensure the radio stops when the main process exits
	defer m.Stop()

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
