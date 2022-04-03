package command

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/aligator/HideAndShell/game/shell"
	"github.com/aligator/HideAndShell/game/shell/command/virus"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/pflag"
)

type Rm struct{}

var (
	ErrSystemFileRemoved = errors.New("system file removed")
)

func (r Rm) Exec(ctx shell.Context, input string) (shell.Context, string, tea.Cmd, error) {
	var flagSet = pflag.NewFlagSet(os.Args[0], pflag.ContinueOnError)
	err := flagSet.Parse(strings.Split(input, " "))
	if err != nil {
		return ctx, "", nil, err
	}

	var args = filterArgs(flagSet.Args())

	if len(args) > 1 || len(args) < 1 {
		return ctx, "", nil, errors.New("too many arguments, 'rm' only supports 1 argument")
	}

	var path string
	if strings.HasPrefix(args[0], "/") {
		path = flagSet.Arg(0)
	} else {
		path = filepath.Join(ctx.WorkingDirectory, args[0])
	}

	file, err := ctx.Filesystem.Open(path)
	if err != nil {
		return ctx, "", nil, err
	}
	var startOfFile = make([]byte, 10)
	_, err = file.Read(startOfFile)
	if err != nil {
		return ctx, "", nil, err
	}
	file.Close()
	if bytes.HasPrefix(startOfFile, []byte("sys")) {
		return ctx, "", nil, ErrSystemFileRemoved
	}

	if bytes.HasPrefix(startOfFile, []byte("virus")) {
		if ctx.Virus.PID != 0 {
			return ctx, "", nil, errors.New("cannot delete file - it is currently used")
		}
		return ctx, "", virus.DeletedCmd, nil
	}

	err = ctx.Filesystem.Remove(path)
	if err != nil {
		return ctx, "", nil, err
	}

	return ctx, "", nil, nil
}
