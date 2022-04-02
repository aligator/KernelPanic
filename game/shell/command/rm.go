package command

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/aligator/HideAndShell/game/shell"
	"github.com/spf13/pflag"
)

type Rm struct{}

func (r Rm) Exec(ctx shell.Context, input string) (shell.Context, string, error) {
	var flagSet = pflag.NewFlagSet(os.Args[0], pflag.ContinueOnError)
	err := flagSet.Parse(strings.Split(input, " "))
	if err != nil {
		return ctx, "", err
	}

	var args = filterArgs(flagSet.Args())

	if len(args) > 1 || len(args) < 1 {
		return ctx, "", errors.New("too many arguments, 'mkdir' only supports 1 argument")
	}

	var path string
	if strings.HasPrefix(args[0], "/") {
		path = flagSet.Arg(0)
	} else {
		path = filepath.Join(ctx.WorkingDirectory, args[0])
	}

	err = ctx.Filesystem.Remove(path)
	if err != nil {
		return ctx, "", err
	}

	return ctx, "", nil
}
