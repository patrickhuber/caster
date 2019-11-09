package vfs

import "github.com/spf13/afero"

func NewMemory() FileSystem {
	fs := afero.NewMemMapFs()
	return NewAfero(fs)
}
