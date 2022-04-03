package virus

import (
	"math/rand"
	"path/filepath"
	"time"

	"github.com/aligator/HideAndShell/game/score"
	"github.com/aligator/HideAndShell/game/shell/filesystem"
	tea "github.com/charmbracelet/bubbletea"
)

var evilFilenames = []string{"evil.exe", "kill-all.exe", "boom.exe", "1337.exe", "LOVE-LETTER-FOR-YOU.TXT.vbs", "NUKE.exe"}

type tickDisableAlertMsg struct{}
type enableAlertMsg struct{}

var alertCmd = tea.Batch(tea.Tick(time.Second*5, func(t time.Time) tea.Msg {
	return tickDisableAlertMsg{}
}), func() tea.Msg {
	return enableAlertMsg{}
})

var Viruses = []Virus{
	BombModel{},
}

type tickMsg time.Time

func KilledCmd() tea.Msg {
	return killedMsg{}
}

func DeletedCmd() tea.Msg {
	return deletedMsg{}
}

func restartVirusCmd() tea.Msg {
	return restartVirusMsg{}
}

type killedMsg struct{}
type deletedMsg struct{}
type restartVirusMsg struct{}

type Virus interface {
	Update(msg tea.Msg) (Virus, tea.Cmd)
	Init() tea.Cmd
	View() string
}

type Model struct {
	Filesystem      *filesystem.Filesystem
	PID             int
	currentlyActive Virus
	CurrentLocation string
	CurrentName     string
}

var tickCmd = tea.Tick(time.Second, func(t time.Time) tea.Msg {
	return tickMsg(t)
})

func (m Model) Init() tea.Cmd {
	return tickCmd
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg.(type) {
	case tickMsg:
		if m.currentlyActive == nil && rand.Intn(3) == 1 {

			// Spawn new virus
			m.currentlyActive = Viruses[rand.Intn(len(Viruses))]

			folderPaths := make([]string, len(m.Filesystem.AllFolders))

			i := 0
			for k := range m.Filesystem.AllFolders {
				folderPaths[i] = k
				i++
			}

			randomFolder := folderPaths[rand.Intn(len(folderPaths))]

			evilFilename := evilFilenames[rand.Intn(len(evilFilenames))]

			file, err := m.Filesystem.Create(filepath.Join(randomFolder, evilFilename))
			if err != nil {
				panic(err)
			}
			_, err = file.Write([]byte("virus"))
			if err != nil {
				panic(err)
			}
			file.Close()
			m.CurrentLocation = filepath.Join(randomFolder, evilFilename)
			m.CurrentName = evilFilename
			m.PID = rand.Intn(999) + 1000

			return m, m.currentlyActive.Init()
		}

		cmds = append(cmds, tickCmd)
	case deletedMsg:
		m.currentlyActive = nil
		m.PID = 0
		m.CurrentName = ""
		m.CurrentLocation = ""

		return m, score.AddScoreCmd(1)
	case restartVirusMsg:
		if m.currentlyActive != nil {
			m.PID = rand.Intn(999) + 1000
		}
	}

	if m.currentlyActive != nil {
		m.currentlyActive, cmd = m.currentlyActive.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.currentlyActive == nil {
		return "No virus found..."
	}

	return m.currentlyActive.View()
}
