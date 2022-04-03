package command

import (
	"github.com/aligator/HideAndShell/game/shell"
	tea "github.com/charmbracelet/bubbletea"
)

type Command interface {
	Exec(ctx shell.Context, input string) (shell.Context, string, tea.Cmd, error)
}
