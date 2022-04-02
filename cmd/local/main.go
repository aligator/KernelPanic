package main

import (
	"log"

	"github.com/aligator/HideAndShell/game"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	g := game.New()
	p := tea.NewProgram(g, tea.WithMouseAllMotion())

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
