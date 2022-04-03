package command

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/aligator/HideAndShell/game/shell"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/pflag"
)

type Mkdir struct{}

func (m Mkdir) Exec(ctx shell.Context, input string) (shell.Context, string, tea.Cmd, error) {
	var flagSet = pflag.NewFlagSet(os.Args[0], pflag.ContinueOnError)
	err := flagSet.Parse(strings.Split(input, " "))
	if err != nil {
		return ctx, "", nil, err
	}

	var args = filterArgs(flagSet.Args())

	if len(args) > 1 || len(args) < 1 {
		return ctx, "", nil, errors.New("too many arguments, 'mkdir' only supports 1 argument")
	}

	var path string
	if strings.HasPrefix(args[0], "/") {
		path = flagSet.Arg(0)
	} else {
		path = filepath.Join(ctx.WorkingDirectory, args[0])
	}

	err = ctx.Filesystem.Mkdir(path, os.ModeDir)
	if err != nil {
		return ctx, "", nil, err
	}

	return ctx, "", nil, nil
}
