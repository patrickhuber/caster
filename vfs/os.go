package vfs

import "github.com/spf13/afero"

// NewOs creates an os based FileSystem
func NewOs() FileSystem {
	fs := afero.NewOsFs()
	return NewAfero(fs)
}
