package vfs

import (
	"io"
	"os"
)

// FileSystem provides an abstract interface for file system operations
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
}
