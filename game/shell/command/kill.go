package command

import (
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/aligator/HideAndShell/game/shell"
	"github.com/aligator/HideAndShell/game/shell/command/virus"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/pflag"
)

type Kill struct{}

var (
	ErrSystemPIDKilled = errors.New("system pid killed")
)

func (k Kill) Exec(ctx shell.Context, input string) (shell.Context, string, tea.Cmd, error) {
	var flagSet = pflag.NewFlagSet(os.Args[0], pflag.ContinueOnError)
	err := flagSet.Parse(strings.Split(input, " "))
	if err != nil {
		return ctx, "", nil, err
	}

	var args = filterArgs(flagSet.Args())

	if len(args) > 1 || len(args) < 1 {
		return ctx, "", nil, errors.New("too many arguments, 'kill' only supports 1 argument")
	}

	pid, err := strconv.Atoi(args[0])
	if err != nil {
		return ctx, "", nil, errors.New("use the PID and not the Process name")
	}

	if len(randomProcessesPIDs) == 0 {
		initRandomProcesses()
	}
	found := false
	for _, d := range randomProcessesPIDs {
		if d == pid {
			found = true
			break
		}
	}

	if found {
		return ctx, "", nil, ErrSystemPIDKilled
	}

	if ctx.Virus.PID != 0 && pid == ctx.Virus.PID {
		ctx.Virus.PID = 0
		return ctx, "Process killed", virus.KilledCmd, nil
	}

	return ctx, "no process killed", nil, nil
}
