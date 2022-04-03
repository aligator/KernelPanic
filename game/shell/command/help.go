package command

import (
	"github.com/aligator/HideAndShell/game/shell"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Help struct{}

var boldStyle = lipgloss.NewStyle().Bold(true)

func (h Help) Exec(ctx shell.Context, input string) (shell.Context, string, tea.Cmd, error) {
	help := boldStyle.Render("Goal:")
	help += `
You are the administrator of this server.
Some hackers plant viruses.
You have to stop them to avoid a kernel panic!

`

	help += boldStyle.Render("Scrolling:")
	help += `
In some terminals the scroll wheel just works.
In others, just use the "PageUp" and "PageDown" or "ctrl+u" and "ctrl+d"

`
	help += boldStyle.Render("Quit game:")
	help += `
"ctrl+c" or type "exit"

`
	help += boldStyle.Render("Available commands:")
	help += `
• ls {Path}    (list directory)
• mkdir {Path} (create folder)
• cd {Path}    (change directory)
• rm {Path}    (delete file)
• ps           (list processes)
• kill {PID}   (kill process)

See further instructions bellow...
`
	return ctx, help, nil, nil
}
