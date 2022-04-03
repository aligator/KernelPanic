package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/aligator/HideAndShell/game"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	rand.Seed(time.Now().UnixMilli())

	g := game.New()
	p := tea.NewProgram(g, tea.WithMouseAllMotion())

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
