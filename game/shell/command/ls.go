package command

import (
	"errors"
	"os"
	"strings"

	"github.com/aligator/HideAndShell/game/shell"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/afero"
	"github.com/spf13/pflag"
)

var styleFolder = lipgloss.NewStyle().Bold(true)
var styleFile = lipgloss.NewStyle().Bold(false)

type Ls struct{}

func (l Ls) Exec(ctx shell.Context, input string) (shell.Context, string, error) {
	var flagSet = pflag.NewFlagSet(os.Args[0], pflag.ContinueOnError)
	err := flagSet.Parse(strings.Split(input, " "))
	if err != nil {
		return ctx, "", err
	}

	flagSet.Args()

	path := ctx.WorkingDirectory
	if len(flagSet.Args()) > 0 {
		path = flagSet.Arg(0)
	}
	if len(flagSet.Args()) > 1 {
		return ctx, "", errors.New("too many arguments, 'ls' only supports up to 1 arguments")
	}

	stat, err := ctx.Filesystem.Stat(path)
	if err != nil {
		return ctx, "", err
	}

	if !stat.IsDir() {
		return ctx, styleFile.Render(stat.Name()), nil
	}

	dir, err := afero.ReadDir(ctx.Filesystem, path)
	if err != nil {
		return ctx, "", err
	}

	output := ""

	for _, file := range dir {
		if file.IsDir() {
			output += styleFolder.Render(file.Name()) + "\n"
		} else {
			output += styleFile.Render(file.Name()) + "\n"
		}
	}

	output = strings.TrimSuffix(output, "\n")

	return ctx, output, nil
}
