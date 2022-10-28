package fs

import "github.com/spf13/afero"

// NewOs creates an os based FileSystem
func NewOs() FS {
	fs := afero.NewOsFs()
	return NewAfero(fs)
}
