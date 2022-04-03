package bluescreen

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	active bool
}

var bsodStyle = lipgloss.NewStyle().Background(lipgloss.Color("#0000ff")).Foreground(lipgloss.Color("#ffffff")).Bold(true)

func NewBlueScreeneModel() Model {
	m := Model{}
	return m
}

type activateMsg struct{}

var BSODCmd = func() tea.Msg {
	return activateMsg{}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		bsodStyle.Width(msg.Width)
		bsodStyle.Height(msg.Height)
	case activateMsg:
		m.active = true
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.active {
		return bsodStyle.Render("YOU ARE DEAD!!!")
	}

	return ""
}
