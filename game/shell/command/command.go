package command

import "github.com/aligator/HideAndShell/game/shell"

type Command interface {
	Exec(ctx shell.Context, input string) (shell.Context, string, error)
}
