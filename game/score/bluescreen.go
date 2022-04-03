package score

import (
	"strconv"
	"strings"
	"sync"

	"github.com/aligator/HideAndShell/server"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type BluescreenModel struct {
	reason    string
	score     int
	bsodStyle lipgloss.Style

	nameInput textinput.Model

	serverHighscore *server.Highscore

	highscoreSnapshot []server.Score

	submitMutex sync.Mutex

	highscoreViewport viewport.Model

	ready bool
}

func NewBlueScreenModel(serverHighscore *server.Highscore) BluescreenModel {
	nameInput := textinput.New()
	nameInput.Focus()
	nameInput.CharLimit = 10
	nameInput.Prompt = "Enter your name:"

	m := BluescreenModel{
		serverHighscore: serverHighscore, reason: "some reason",
		bsodStyle: lipgloss.NewStyle().Background(lipgloss.Color("#0000ff")).Foreground(lipgloss.Color("#ffffff")).Bold(true),
		nameInput: nameInput,
	}

	return m
}

type activateMsg struct {
	reason string
}

func BSODCmd(reason string) tea.Cmd {
	return func() tea.Msg {
		return activateMsg{
			reason,
		}
	}
}

func (m BluescreenModel) Init() tea.Cmd {
	return nil
}

func (m BluescreenModel) Update(msg tea.Msg) (BluescreenModel, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.bsodStyle.Width(msg.Width)
		m.bsodStyle.Height(msg.Height)

		if !m.ready {
			m.highscoreViewport = viewport.New(0, msg.Height-9-13)
			m.highscoreViewport.SetContent("")
			m.highscoreViewport.GotoTop()
			m.highscoreViewport.KeyMap = viewport.DefaultKeyMap()
			m.highscoreViewport.KeyMap.HalfPageUp = key.NewBinding(
				key.WithKeys("ctrl+u"),
			)
			m.highscoreViewport.KeyMap.HalfPageDown = key.NewBinding(
				key.WithKeys("ctrl+d"),
			)
			m.highscoreViewport.KeyMap.Up = key.NewBinding(
				key.WithKeys("up"),
			)
			m.highscoreViewport.KeyMap.Down = key.NewBinding(
				key.WithKeys("down"),
			)

			m.ready = true
		} else {
			m.highscoreViewport.Height = msg.Height - 9 - 12
		}
	case activateMsg:
		m.reason = msg.reason
	case scoreMsg:
		m.score += int(msg)
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.highscoreSnapshot == nil && m.reason != "" {
				m.serverHighscore.Insert(server.Score{
					Name:  m.nameInput.Value(),
					Value: m.score,
				})
				m.highscoreSnapshot = m.serverHighscore.Get()

				highscore := ""
				for i, score := range m.highscoreSnapshot {
					highscore += strconv.Itoa(i+1) + ": " + strconv.Itoa(score.Value) + " - " + score.Name + "\n"
				}

				highscore = strings.TrimSuffix(highscore, "\n")

				m.highscoreViewport.SetContent(highscore)
				m.highscoreViewport.Width = lipgloss.Width(highscore)
				m.highscoreViewport.GotoTop()

			}
		}
	}

	if m.reason != "" {
		m.nameInput, cmd = m.nameInput.Update(msg)
		cmds = append(cmds, cmd)

		m.highscoreViewport, cmd = m.highscoreViewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m BluescreenModel) View() string {
	if m.reason == "" {
		return ""
	}

	kernelPanic := "KERNEL PANIC"
	reason := m.reason
	score := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Render("Your Score: " + strconv.Itoa(m.score))
	highscore := ""
	if m.highscoreSnapshot == nil {
		highscore += m.nameInput.View()
	} else {
		highscore += "Highscore:\n"
	}

	header := lipgloss.JoinVertical(lipgloss.Center, kernelPanic, reason, score, `--------------
/                \
/  ---        ---  \
| |   |      |   | |
|  ---        ---  |
\        ^         /
\      / \       /
\_   /___\    _/
| |========| |
\  |======|  /
\          /
--------

`, highscore)

	headerWidth := lipgloss.Width(header)
	header = lipgloss.NewStyle().
		MarginLeft(m.bsodStyle.GetWidth()/2 - headerWidth/2).
		Render(header)

	msg := header + "\n"

	if m.highscoreSnapshot != nil {
		board := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1).Render(m.highscoreViewport.View())

		board = lipgloss.NewStyle().
			MarginLeft(m.bsodStyle.GetWidth()/2 - lipgloss.Width(board)/2).
			Render(board)

		msg += board
	}

	return m.bsodStyle.Render(msg)
}
