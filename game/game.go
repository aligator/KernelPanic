package game

import (
	"errors"
	"fmt"
	"strings"

	"github.com/aligator/HideAndShell/game/bluescreen"
	"github.com/aligator/HideAndShell/game/score"
	"github.com/aligator/HideAndShell/game/shell"
	"github.com/aligator/HideAndShell/game/shell/command"
	"github.com/aligator/HideAndShell/game/shell/command/virus"
	"github.com/aligator/HideAndShell/game/shell/filesystem"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type game struct {
	history  shell.HistoryModel
	cmdInput textinput.Model
	score    score.Model
	bsod     bluescreen.Model

	ctx shell.Context

	commands map[string]command.Command
}

var headerStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder())
var virusStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder())
var historyStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder())
var inputStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder())

func New() game {
	cmdInput := textinput.New()
	cmdInput.Focus()

	fs := filesystem.New()

	m := game{
		cmdInput: cmdInput,
		history:  shell.NewHistory(),
		score:    score.NewScoreModel(),
		bsod:     bluescreen.NewBlueScreeneModel(),

		ctx: shell.Context{
			Filesystem:       fs,
			WorkingDirectory: "/",

			Virus: virus.Model{
				Filesystem: fs,
			},
		},
		commands: map[string]command.Command{
			"ls":    command.Ls{},
			"mkdir": command.Mkdir{},
			"cd":    command.Cd{},
			"rm":    command.Rm{},
			"ps":    command.Ps{},
			"kill":  command.Kill{},
		},
	}

	return m
}

func (m game) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.history.Init(), m.ctx.Virus.Init())
}

func (m game) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		headerStyle.Width(msg.Width - headerStyle.GetBorderRightSize() - 1)
		virusStyle.Width(msg.Width - virusStyle.GetBorderRightSize() - 1)
		historyStyle.Width(msg.Width - historyStyle.GetBorderRightSize() - 1)
		inputStyle.Width(msg.Width - inputStyle.GetBorderRightSize() - 1)
		m.history.Top = historyStyle.GetBorderTopWidth() + lipgloss.Height(headerStyle.Render(m.score.View()))
		m.history.Right = historyStyle.GetBorderRightSize()
		m.history.Bottom = historyStyle.GetBorderBottomSize() + 4
		m.history.Left = historyStyle.GetBorderLeftSize()
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
					m.history, cmd = m.history.Update(shell.AddHistoryMsg{Text: "> " + input + "\n" + "unknown command"})
					return
				}

				shellCmd := m.commands[splitted[0]]

				var result string
				var err error
				m.ctx, result, cmd, err = shellCmd.Exec(m.ctx, splitted[1])
				if err != nil {
					if errors.Is(err, command.ErrSystemPIDKilled) || errors.Is(err, command.ErrSystemFileRemoved) {
						cmds = append(cmds, bluescreen.BSODCmd)
						return
					}

					m.history, cmd = m.history.Update(shell.AddHistoryMsg{Text: "> " + input + "\n" + err.Error()})
					return
				}
				if cmd != nil {
					cmds = append(cmds, cmd)
					cmd = nil
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

	m.ctx.Virus, cmd = m.ctx.Virus.Update(msg)
	cmds = append(cmds, cmd)

	m.score, cmd = m.score.Update(msg)
	cmds = append(cmds, cmd)

	m.bsod, cmd = m.bsod.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m game) View() string {
	bsod := m.bsod.View()
	if bsod != "" {
		return bsod
	}
	return fmt.Sprintf("%s\n%s\n%s\n%s", headerStyle.Render(m.score.View()), historyStyle.Render(m.history.View()), virusStyle.Render(m.ctx.Virus.View()), m.cmdInput.View())
}
