package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/aligator/HideAndShell/game"
	"github.com/aligator/HideAndShell/server"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	rand.Seed(time.Now().UnixMilli())

	highscore := &server.Highscore{}

	g := game.New(highscore)
	p := tea.NewProgram(g, tea.WithMouseAllMotion())

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
