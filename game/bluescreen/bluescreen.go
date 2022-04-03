package bluescreen

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	active    bool
	bsodStyle lipgloss.Style
}

func NewBlueScreeneModel() Model {
	m := Model{
		bsodStyle: lipgloss.NewStyle().Background(lipgloss.Color("#0000ff")).Foreground(lipgloss.Color("#ffffff")).Bold(true),
	}
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
		m.bsodStyle.Width(msg.Width)
		m.bsodStyle.Height(msg.Height)
	case activateMsg:
		m.active = true
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.active {
		return m.bsodStyle.Render("YOU ARE DEAD!!!")
	}

	return ""
}
