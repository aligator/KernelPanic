package score

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	score int
}

var scoreStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Bold(true)

func NewScoreModel() Model {
	m := Model{
		score: 0,
	}

	return m
}

type scoreMsg int

func AddScoreCmd(count int) tea.Cmd {
	return func() tea.Msg {
		return scoreMsg(count)
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case scoreMsg:
		m.score += int(msg)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return fmt.Sprintf(scoreStyle.Render("Score: %v"), m.score)
}
