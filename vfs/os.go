package vfs

import "github.com/spf13/afero"

func NewOs() FileSystem {
	fs := afero.NewOsFs()
	return NewAfero(fs)
}
