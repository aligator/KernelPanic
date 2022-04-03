package score

import (
	"strconv"
	"sync"

	"github.com/aligator/HideAndShell/server"
	"github.com/charmbracelet/bubbles/textinput"
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
}

func NewBlueScreenModel(serverHighscore *server.Highscore) BluescreenModel {
	nameInput := textinput.New()
	nameInput.Focus()
	nameInput.CharLimit = 10
	nameInput.Prompt = "Enter your name:"

	m := BluescreenModel{
		serverHighscore: serverHighscore,
		/*reason:          "some reason",
		score:       0,*/
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
	case activateMsg:
		m.reason = msg.reason
	case scoreMsg:
		m.score += int(msg)
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}

		switch msg.Type {
		case tea.KeyEnter:
			if m.highscoreSnapshot == nil {
				m.serverHighscore.Insert(server.Score{
					Name:  m.nameInput.Value(),
					Value: m.score,
				})
				m.highscoreSnapshot = m.serverHighscore.Get()
			}
		}
	}

	m.nameInput, cmd = m.nameInput.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m BluescreenModel) View() string {
	if m.reason != "" {
		msg := "YOU ARE DEAD!!!\n" + m.reason + "\nYour Score: " + strconv.Itoa(m.score)
		msg += "\n\n"

		if m.highscoreSnapshot == nil {
			msg += m.nameInput.View() + "\n\n"
		} else {
			msg += "Highscore:"
			for i, score := range m.highscoreSnapshot {
				msg += strconv.Itoa(i) + ": " + strconv.Itoa(score.Value) + " - " + score.Name + "\n"
			}
		}
		return m.bsodStyle.Render(msg)
	}

	return ""
}
