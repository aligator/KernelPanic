package command

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/aligator/HideAndShell/game/shell"
	"github.com/spf13/pflag"
)

type Cd struct{}

func (c Cd) Exec(ctx shell.Context, input string) (shell.Context, string, error) {
	var flagSet = pflag.NewFlagSet(os.Args[0], pflag.ContinueOnError)
	err := flagSet.Parse(strings.Split(input, " "))
	if err != nil {
		return ctx, "", err
	}

	var args = filterArgs(flagSet.Args())

	if len(args) > 1 || len(args) < 1 {
		return ctx, "", errors.New("too many arguments, 'cd' only supports 1 argument")
	}

	var path string
	if strings.HasPrefix(args[0], "/") {
		path = args[0]
	} else {
		path, err = filepath.Abs(filepath.Join(ctx.WorkingDirectory, args[0]))
		if err != nil {
			return ctx, "", err
		}
	}

	stat, err := ctx.Filesystem.Stat(path)
	if err != nil {
		return ctx, "", err
	}

	if !stat.IsDir() {
		return ctx, "", errors.New("you can only cd into a directory, not into a file")
	}

	ctx.WorkingDirectory = path

	return ctx, "", nil
}
