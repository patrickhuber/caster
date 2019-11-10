package vfs

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

// FileSystem provides an abstract interface for file system operations. It's main purpose is to prevent leaking provider apis across the project.
type FileSystem interface {
	Rename(oldname, newname string) error
	Create(path string) (File, error)
	Write(path string, data []byte, permissions os.FileMode) error
	Exists(path string) (bool, error)
	IsDir(path string) (bool, error)
	Mkdir(path string, permissions os.FileMode) error
	Stat(name string) (os.FileInfo, error)
	Open(name string) (File, error)
	WriteReader(path string, reader io.Reader) error
	RemoveAll(path string) error
	Remove(name string) error
	Read(filename string) ([]byte, error)
	ReadDir(dirname string) ([]os.FileInfo, error)
	Walk(root string, walkFn filepath.WalkFunc) error
	Join(segments ...string) string
}

// fileSystem provides default implementaions for functions that are cross provider
type fileSystem struct {
}

func (fs *fileSystem) Join(segments ...string) string {
	path := strings.Join(segments, "/")
	if len(segments) == 1 {
		if segments[0] == "" {
			path = "/" + path
		}
	}
	return fs.ToSlash(path)
}

// ToSlash replaces all backslashes with forward slashes
func (fs *fileSystem) ToSlash(path string) string {
	return strings.Replace(path, "\\", "/", -1)
}
