package vfs

import (
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

type aferoFs struct {
	*fileSystem
	fs afero.Fs
}

// NewAfero creates a FileSystem implemented with afero.Fs
func NewAfero(fs afero.Fs) FileSystem {
	return &aferoFs{
		fs: fs,
	}
}

func (fs *aferoFs) Rename(oldname, newname string) error {
	return fs.fs.Rename(oldname, newname)
}

func (fs *aferoFs) Create(path string) (File, error) {
	return fs.fs.Create(path)
}

func (fs *aferoFs) Write(path string, data []byte, permissions os.FileMode) error {
	return afero.WriteFile(fs.fs, path, data, permissions)
}

func (fs *aferoFs) Exists(path string) (bool, error) {
	return afero.Exists(fs.fs, path)
}

func (fs *aferoFs) IsDir(path string) (bool, error) {
	return afero.IsDir(fs.fs, path)
}

func (fs *aferoFs) Mkdir(path string, permissions os.FileMode) error {
	return fs.fs.Mkdir(path, permissions)
}

func (fs *aferoFs) Stat(name string) (os.FileInfo, error) {
	return fs.fs.Stat(name)
}

func (fs *aferoFs) Open(name string) (File, error) {
	return fs.fs.Open(name)
}

func (fs *aferoFs) WriteReader(path string, reader io.Reader) error {
	return afero.WriteReader(fs.fs, path, reader)
}

func (fs *aferoFs) RemoveAll(path string) error {
	return fs.fs.RemoveAll(path)
}

func (fs *aferoFs) Remove(name string) error {
	return fs.fs.Remove(name)
}

func (fs *aferoFs) Read(filename string) ([]byte, error) {
	return afero.ReadFile(fs.fs, filename)
}

func (fs *aferoFs) ReadDir(dirname string) ([]os.FileInfo, error) {
	return afero.ReadDir(fs.fs, dirname)
}

func (fs *aferoFs) Walk(root string, walkFn filepath.WalkFunc) error {
	return afero.Walk(
		fs.fs,
		root,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			path = fs.ToSlash(path)
			return walkFn(path, info, err)
		})
}
