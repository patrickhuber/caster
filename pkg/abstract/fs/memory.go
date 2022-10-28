package fs

import "github.com/spf13/afero"

// NewMemory creates a new in memory FileSystem
func NewMemory() FS {
	fileSystem := afero.NewMemMapFs()
	return NewAfero(fileSystem)
}
