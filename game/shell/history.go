package shell

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type HistoryModel struct {
	lines           []string
	historyViewport viewport.Model
	footerHeight    int

	ready bool
}

func NewHistory(footerHeight int) HistoryModel {
	m := HistoryModel{
		lines:        strings.Split(``, "\n"),
		footerHeight: footerHeight,
	}

	return m
}

func (m HistoryModel) String() string {
	return strings.Join(m.lines, "\n")
}

type AddHistoryMsg struct {
	Text string
}

func (m *HistoryModel) add(line string) {
	m.lines = append(m.lines, strings.Split(line, "\n")...)
	m.historyViewport.SetContent(m.String())
	m.historyViewport.GotoBottom()
}

func (m HistoryModel) Init() tea.Cmd {
	return nil
}

func (m HistoryModel) Update(msg tea.Msg) (HistoryModel, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		if !m.ready {
			m.historyViewport = viewport.New(msg.Width, msg.Height-m.footerHeight)
			m.historyViewport.SetContent(m.String())
			m.historyViewport.GotoBottom()
			m.historyViewport.KeyMap = viewport.DefaultKeyMap()
			m.historyViewport.KeyMap.HalfPageUp = key.NewBinding(
				key.WithKeys("ctrl+u"),
			)
			m.historyViewport.KeyMap.HalfPageDown = key.NewBinding(
				key.WithKeys("ctrl+d"),
			)
			m.historyViewport.KeyMap.Up = key.NewBinding(
				key.WithKeys("up"),
			)
			m.historyViewport.KeyMap.Down = key.NewBinding(
				key.WithKeys("down"),
			)

			m.ready = true
		} else {
			m.historyViewport.Width = msg.Width
			m.historyViewport.Height = msg.Height - m.footerHeight
		}
	case AddHistoryMsg:
		m.add(msg.Text)
	}

	m.historyViewport, cmd = m.historyViewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m HistoryModel) View() string {
	if !m.ready {
		return "LOADING..."
	}
	return m.historyViewport.View()
}
