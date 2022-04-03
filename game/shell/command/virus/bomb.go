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
}

var tickBombCmd = tea.Tick(time.Second, func(t time.Time) tea.Msg {
	return tickBombMsg(t)
})

func (m BombModel) Init() tea.Cmd {
	return tea.Batch(alertCmd, tickBombCmd)
}

type tickBombMsg time.Time
type tickDisableAlertMsg struct{}
type enableAlertMsg struct{}

var alertCmd = tea.Batch(tea.Tick(time.Second*5, func(t time.Time) tea.Msg {
	return tickDisableAlertMsg{}
}), func() tea.Msg {
	return enableAlertMsg{}
})

func (m BombModel) Update(msg tea.Msg) (Virus, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg.(type) {
	case restartVirusMsg:
		m.endTime = time.Now().Add(time.Second * (time.Duration(rand.Intn(50) + 30)))
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
				cmds = append(cmds, restartVirusCmd, alertCmd)
			} else {
				cmds = append(cmds, score.BSODCmd("'bomb!.worm' destroyed your system"))
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
		msg += "Successfully stopped virus! Find and delete the 'bomb!.worm' or it may be restarted!!!"
	} else {
		msg += "Virus of type 'bomb!.worm' found! Kill the process quickly!!! (hint: 'ps' & 'kill')"
	}

	if !m.endTime.IsZero() {
		msg = strconv.FormatFloat(m.endTime.Sub(time.Now()).Seconds(), 'f', 0, 64) + " " + msg
	}

	return m.styleAlert.Render(msg)
}
