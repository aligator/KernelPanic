package virus

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/aligator/HideAndShell/game/score"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type BombModel struct {
	killed     bool
	endTime    time.Time
	styleAlert lipgloss.Style
	stopped    bool
}

var tickBombCmd = tea.Tick(time.Second, func(t time.Time) tea.Msg {
	return tickBombMsg(t)
})

func (m BombModel) Init() tea.Cmd {
	return tea.Batch(alertCmd, tickBombCmd)
}

type tickBombMsg time.Time

func (m BombModel) Update(msg tea.Msg) (Virus, tea.Cmd) {
	var cmds []tea.Cmd

	if m.stopped {
		return m, nil
	}

	switch msg.(type) {
	case restartVirusMsg:
		m.endTime = time.Now().Add(time.Second * (time.Duration(rand.Intn(30) + 20)))
	case killedMsg:
		m.killed = true
		m.endTime = time.Now().Add(time.Second * time.Duration(rand.Intn(20)+60))
	case tickBombMsg:
		if m.endTime.IsZero() {
			m.endTime = time.Now().Add(time.Second * (time.Duration(rand.Intn(50) + 30)))
		}

		if time.Now().UnixMilli() >= m.endTime.UnixMilli() {
			if m.killed {
				m.killed = false
				return m, tea.Batch(restartVirusCmd, alertCmd)
			} else {
				m.stopped = true
				return m, score.BSODCmd("'bomb!.worm' destroyed your system")
			}
		}
		cmds = append(cmds, tickBombCmd)

	case tickDisableAlertMsg:
		m.styleAlert = lipgloss.NewStyle().Background(lipgloss.Color("0")).Foreground(lipgloss.Color("2"))
	case enableAlertMsg:
		m.styleAlert = lipgloss.NewStyle().Background(lipgloss.Color("0")).Foreground(lipgloss.Color("2")).Blink(true)
	}

	return m, tea.Batch(cmds...)
}

func (m BombModel) View() string {
	msg := ""

	if m.killed {
		msg += "Successfully stopped virus! Find and delete executable file or it may be restarted!!! (hint: 'ls', 'cd', 'rm')"
	} else {
		msg += "Virus found! Kill the process quickly!!! (hint: 'ps' & 'kill')"
	}

	if !m.endTime.IsZero() {
		msg = strconv.FormatFloat(m.endTime.Sub(time.Now()).Seconds(), 'f', 0, 64) + " " + msg
	}

	return m.styleAlert.Render(msg)
}
