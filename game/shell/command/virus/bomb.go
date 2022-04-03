package virus

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/aligator/HideAndShell/game/bluescreen"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type BombModel struct {
	killed  bool
	alert   bool
	endTime time.Time
}

var styleAlert = lipgloss.NewStyle().Blink(true)

var tickBombCmd = tea.Tick(time.Second, func(t time.Time) tea.Msg {
	return tickBombMsg(t)
})

func (m BombModel) Init() tea.Cmd {
	return tea.Batch(alertCmd, tickBombCmd)
}

type tickBombMsg time.Time
type tickDisableAlertMsg struct{}
type enableAlertMsg struct{}

var alertCmd = tea.Batch(tea.Tick(time.Second*10, func(t time.Time) tea.Msg {
	return tickDisableAlertMsg{}
}), func() tea.Msg {
	return enableAlertMsg{}
})

func (m BombModel) Update(msg tea.Msg) (Virus, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg.(type) {
	case restartVirusMsg:
		m.endTime = time.Now().Add(time.Second * (time.Duration(rand.Intn(20) + 30)))
	case killedMsg:
		m.killed = true
		m.endTime = time.Now().Add(time.Second * time.Duration(rand.Intn(20)+60))
	case tickBombMsg:
		if m.endTime.IsZero() {
			m.endTime = time.Now().Add(time.Second * (time.Duration(rand.Intn(20) + 30)))
		}

		if time.Now().UnixMilli() >= m.endTime.UnixMilli() {
			if m.killed {
				m.killed = false
				cmds = append(cmds, restartVirusCmd, alertCmd)
			} else {
				cmds = append(cmds, bluescreen.BSODCmd)
			}
		}
		cmds = append(cmds, tickBombCmd)

	case tickDisableAlertMsg:
		m.alert = false
	case enableAlertMsg:
		m.alert = true
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

	if m.alert {
		msg = styleAlert.Render(msg)
	}

	if !m.endTime.IsZero() {
		msg = strconv.FormatFloat(m.endTime.Sub(time.Now()).Seconds(), 'f', 0, 64) + " " + msg
	}

	return msg
}