package command

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/aligator/HideAndShell/game/shell"
	"github.com/spf13/pflag"
)

type Mkdir struct{}

func (m Mkdir) Exec(ctx shell.Context, input string) (shell.Context, string, error) {
	var flagSet = pflag.NewFlagSet(os.Args[0], pflag.ContinueOnError)
	err := flagSet.Parse(strings.Split(input, " "))
	if err != nil {
		return ctx, "", err
	}

	flagSet.Args()

	if len(flagSet.Args()) > 1 || len(flagSet.Args()) < 1 {
		return ctx, "", errors.New("too many arguments, 'mkdir' only supports 1 argument")
	}

	var path string
	if strings.HasPrefix(flagSet.Arg(0), "/") {
		path = flagSet.Arg(0)
	} else {
		path = filepath.Join(ctx.WorkingDirectory, flagSet.Arg(0))
	}

	err = ctx.Filesystem.Mkdir(path, os.ModeDir)
	if err != nil {
		return ctx, "", err
	}

	return ctx, "", nil
}
