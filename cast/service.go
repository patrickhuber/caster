package cast

import (
	"bytes"
	"html"
	"os"
	"strings"
	"text/template"

	"github.com/patrickhuber/caster/vfs"
)

// Service handles casting of a template
type Service interface {
	Cast(source string, target string, data map[string]interface{}) error
}

type service struct {
	fs vfs.FileSystem
}

// NewService creates a new instance of the cast service
func NewService(fs vfs.FileSystem) Service {
	return &service{
		fs: fs,
	}
}

func (s *service) Cast(source, target string, data map[string]interface{}) error {
	return s.fs.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return s.castDirectory(source, target, path, info, err)
		}
		return s.castFile(source, target, path, info, err)
	})
}

func (s *service) castDirectory(source, target, path string, info os.FileInfo, err error) error {
	rel, err := s.fs.Rel(source, path)
	if err != nil {
		return err
	}
	targetpath := s.fs.Join(target, rel, info.Name())
	return s.fs.Mkdir(targetpath, 600)
}

func (s *service) castFile(source, target, path string, info os.FileInfo, err error) error {

	rel, err := s.fs.Rel(source, s.fs.Dir(path))
	if err != nil {
		return err
	}

	if info.Name() == ".caster" {
		return nil
	}

	if !strings.Contains(info.Name(), "{{") {
		targetpath := s.fs.Join(target, rel, info.Name())
		targetpath = s.fs.Clean(targetpath)
		content, err := s.fs.Read(path)
		if err != nil {
			return err
		}
		return s.fs.Write(targetpath, content, 0600)
	}

	name := html.UnescapeString(info.Name())
	t, err := template.New("caster").Parse(name)
	if err != nil {
		return err
	}

	var buffer bytes.Buffer
	err = t.Execute(&buffer, nil)
	if err != nil {
		return err
	}

	// get the directory
	targetpath := s.fs.Join(target, rel, buffer.String())
	targetpath = s.fs.Clean(targetpath)
	content, err := s.fs.Read(path)
	if err != nil {
		return err
	}
	return s.fs.Write(targetpath, content, 0600)
}
