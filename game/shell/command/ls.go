package command

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/aligator/HideAndShell/game/shell"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/afero"
	"github.com/spf13/pflag"
)

var styleFolder = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
var styleFile = lipgloss.NewStyle().Bold(false)

type Ls struct{}

func (l Ls) Exec(ctx shell.Context, input string) (shell.Context, string, tea.Cmd, error) {
	var flagSet = pflag.NewFlagSet(os.Args[0], pflag.ContinueOnError)
	err := flagSet.Parse(strings.Split(input, " "))
	if err != nil {
		return ctx, "", nil, err
	}

	var args = filterArgs(flagSet.Args())

	path := ctx.WorkingDirectory
	if len(args) > 0 {
		if strings.HasPrefix(args[0], "/") {
			path = args[0]
		} else {
			path = filepath.Join(ctx.WorkingDirectory, args[0])
		}
	}
	if len(args) > 1 {
		return ctx, "", nil, errors.New("too many arguments, 'ls' only supports up to 1 arguments")
	}

	stat, err := ctx.Filesystem.Stat(path)
	if err != nil {
		return ctx, "", nil, err
	}

	if !stat.IsDir() {
		return ctx, styleFile.Render(stat.Name()), nil, nil
	}

	dir, err := afero.ReadDir(ctx.Filesystem, path)
	if err != nil {
		return ctx, "", nil, err
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

	return ctx, output, nil, nil
}
