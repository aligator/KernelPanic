package game

import (
	"fmt"

	"github.com/aligator/HideAndShell/game/shell"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type game struct {
	history  shell.HistoryModel
	cmdInput textinput.Model
}

func New() game {
	cmdInput := textinput.New()
	cmdInput.Focus()

	m := game{
		cmdInput: cmdInput,
		history:  shell.NewHistory(1),
	}

	return m
}

func (m game) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.history.Init())
}

func (m game) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}

		switch msg.Type {
		case tea.KeyEnter:
			m.history, cmd = m.history.Update(shell.AddHistoryMsg{Text: m.cmdInput.Value()})
			m.cmdInput.SetValue("")
		}
	}

	m.cmdInput, cmd = m.cmdInput.Update(msg)
	cmds = append(cmds, cmd)

	m.history, cmd = m.history.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m game) View() string {
	return fmt.Sprintf("%s\n%s", m.history.View(), m.cmdInput.View())
}
