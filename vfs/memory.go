package vfs

import "github.com/spf13/afero"

// NewMemory creates a new in memory FileSystem
func NewMemory() FileSystem {
	fs := afero.NewMemMapFs()
	return NewAfero(fs)
}
