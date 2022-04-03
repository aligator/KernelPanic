package filesystem

import (
	"math/rand"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

var casualFilenames = []string{"test.pdf", "Holiday.jpg", "Screenshot.png", "xxx.jpg", "xxx.avi"}

var foldernames = []string{"system", "system32", "os", "etc", "config", "Pictures", "Music", "Documents", "Programming", "Downloads", "holiday"}

type Filesystem struct {
	afero.Fs
	AllFolders map[string]struct{}
}

func (f *Filesystem) Mkdir(name string, perm os.FileMode) error {
	err := f.Fs.Mkdir(name, perm)
	if err != nil {
		return err
	}

	f.AllFolders[name] = struct{}{}
	return nil
}

func (f *Filesystem) MkdirAll(path string, perm os.FileMode) error {
	panic("not implemented")
}

func (f *Filesystem) Remove(name string) error {
	err := f.Fs.Remove(name)
	if err != nil {
		return err
	}

	delete(f.AllFolders, name)
	return nil
}

func (f *Filesystem) RemoveAll(path string) error {
	panic("not implemented")
}

func (f *Filesystem) FolderExists(path string, folderName string) bool {
	_, exists := f.AllFolders[filepath.Join(path, folderName)]
	return exists
}

func New() *Filesystem {
	fs := &Filesystem{
		Fs:         afero.NewMemMapFs(),
		AllFolders: make(map[string]struct{}),
	}

	// Setup some random files and viruses:
	addedFiles := make(map[string]struct{})
	fileExists := func(path string, folderName string) bool {
		_, exists := addedFiles[filepath.Join(path, folderName)]
		return exists
	}

	var folders []string
	var addFolders func(deph int, path string)
	addFolders = func(depth int, path string) {
		if depth > 3 {
			return
		}
		folderCount := rand.Intn(3)
		if depth == 0 {
			folderCount += 2
		}

		for i := 0; i < folderCount; i++ {

			folderName := foldernames[rand.Intn(len(foldernames))]
			for ; fileExists(path, folderName); folderName = foldernames[rand.Intn(len(foldernames))] {
				// again
			}

			newFolderPath := filepath.Join(path, folderName)
			addedFiles[newFolderPath] = struct{}{}

			folders = append(folders, newFolderPath)
			err := fs.Mkdir(newFolderPath, os.ModeDir)
			if err != nil {
				panic(err)
			}

			goDeeper := rand.Intn(2) == 1
			if goDeeper {
				addFolders(depth+1, newFolderPath)
			}
		}
	}

	addFolders(0, "/")

	for _, folder := range append(folders, "/") {
		normalFileCount := rand.Intn(10)
		for i := 0; i < normalFileCount; i++ {
			filename := casualFilenames[rand.Intn(len(casualFilenames))]
			for ; fileExists(folder, filename); filename = casualFilenames[rand.Intn(len(casualFilenames))] {
				// again
			}

			f, err := fs.Create(filepath.Join(folder, filename))
			if err != nil {
				panic(err)
			}
			f.Write([]byte("user"))
			f.Close()
		}

		systemFileCount := rand.Intn(5)
		for i := 0; i < systemFileCount; i++ {
			filename := systemFilenames[rand.Intn(len(systemFilenames))]
			for ; fileExists(folder, filename); filename = systemFilenames[rand.Intn(len(systemFilenames))] {
				// again
			}

			f, err := fs.Create(filepath.Join(folder, filename))
			if err != nil {
				panic(err)
			}
			f.Write([]byte("sys"))
			f.Close()
		}
	}

	return fs
}
