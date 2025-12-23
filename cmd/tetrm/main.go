package main

import (
	"log"

	"termino/internal/game"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(game.NewModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
