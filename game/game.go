package game

import (
	"fmt"
	"strings"

	"github.com/aligator/HideAndShell/game/shell"
	"github.com/aligator/HideAndShell/game/shell/command"
	"github.com/aligator/HideAndShell/game/shell/filesystem"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type game struct {
	history  shell.HistoryModel
	cmdInput textinput.Model

	ctx shell.Context

	commands map[string]command.Command
}

func New() game {
	cmdInput := textinput.New()
	cmdInput.Focus()

	m := game{
		cmdInput: cmdInput,
		history:  shell.NewHistory(1),

		ctx: shell.Context{
			Filesystem:       filesystem.New(),
			WorkingDirectory: "/",
		},
		commands: map[string]command.Command{
			"ls": command.Ls{},
		},
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
			// Evaluate input
			cmd = nil
			func() {
				input := m.cmdInput.Value()
				m.cmdInput.SetValue("")

				splitted := strings.SplitN(input, " ", 2)
				if len(splitted) < 1 {
					return
				} else if len(splitted) < 2 {
					splitted = append(splitted, "")
				}

				if _, ok := m.commands[splitted[0]]; !ok {
					m.history, cmd = m.history.Update(shell.AddHistoryMsg{Text: "unknown command"})
					return
				}

				shellCmd := m.commands[splitted[0]]

				var result string
				var err error
				m.ctx, result, err = shellCmd.Exec(m.ctx, splitted[1])
				if err != nil {
					m.history, cmd = m.history.Update(shell.AddHistoryMsg{Text: "> " + input + "\n" + err.Error()})
					return
				}

				m.history, cmd = m.history.Update(shell.AddHistoryMsg{Text: "> " + input + "\n" + result})
			}()

			if cmd != nil {
				cmds = append(cmds, cmd)
			}
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
