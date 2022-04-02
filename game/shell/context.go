package shell

import "github.com/aligator/HideAndShell/game/shell/filesystem"

type Context struct {
	WorkingDirectory string
	Filesystem       *filesystem.Filesystem
}
