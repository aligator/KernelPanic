package filesystem

import (
	"os"

	"github.com/spf13/afero"
)

type Filesystem struct {
	afero.Fs
}

func New() *Filesystem {
	fs := &Filesystem{
		afero.NewMemMapFs(),
	}

	err := fs.Mkdir("/lol", os.ModeDir)
	if err != nil {
		panic(err)
	}

	return fs
}
