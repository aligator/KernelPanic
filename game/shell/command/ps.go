package command

import (
	"math/rand"
	"sort"
	"strconv"
	"strings"

	"github.com/aligator/HideAndShell/game/shell"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var headerStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))

var processNames = []string{"explorer.exe", "smss.exe", "winit.exe", "services.exe", "antivir.exe", "svchost.exe"}
var randomProcesses []string
var randomProcessesPIDs []int

func initRandomProcesses() {
	for _, name := range processNames {
		pid := rand.Intn(999) + 1000
		randomProcessesPIDs = append(randomProcessesPIDs, pid)
		randomProcesses = append(randomProcesses, strconv.Itoa(pid)+"    "+name)
	}
}

type Ps struct{}

func (p Ps) Exec(ctx shell.Context, input string) (shell.Context, string, tea.Cmd, error) {
	if randomProcesses == nil {
		initRandomProcesses()
	}

	var currentProcessList []string = make([]string, len(randomProcesses))
	copy(currentProcessList, randomProcesses)

	if ctx.Virus.PID != 0 {
		currentProcessList = append(currentProcessList, strconv.Itoa(ctx.Virus.PID)+"    "+ctx.Virus.CurrentName)
	}

	// shuffle
	sort.Slice(currentProcessList, func(i, j int) bool {
		return rand.Intn(2) == 1
	})

	return ctx, headerStyle.Render("PID     Process\n───────────────") + "\n" + strings.Join(currentProcessList, "\n"), nil, nil
}
