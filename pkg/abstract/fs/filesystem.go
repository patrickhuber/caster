package fs

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

// FS provides an abstract interface for file system operations. It's main purpose is to prevent leaking provider apis across the project.
type FS interface {
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
	// Rel returns the relative path
	Rel(basepath, targetpath string) (string, error)
	Dir(path string) string
	Clean(path string) string
	ReadDirRegex(dirname string, filter string) ([]os.FileInfo, error)
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

func (fs *fileSystem) Rel(basepath, targetpath string) (string, error) {
	rel, err := filepath.Rel(basepath, targetpath)
	if err != nil {
		return "", err
	}
	return fs.ToSlash(rel), nil
}

func (fs *fileSystem) Dir(path string) string {
	return fs.ToSlash(filepath.Dir(path))
}

func (fs *fileSystem) Clean(path string) string {
	return fs.ToSlash(filepath.Clean(path))
}
