package initialize

import (
	"errors"
	"io/fs"

	"github.com/patrickhuber/go-xplat/filepath"
	afs "github.com/patrickhuber/go-xplat/fs"
)

type Request struct {
	Template string `yaml:"omitempty"`
}

type Response struct{}

type Service interface {
	Initialize(req *Request) (*Response, error)
}

func NewService(fs afs.FS, path *filepath.Processor) Service {
	return &service{
		fs:   fs,
		path: path,
	}
}

type service struct {
	fs   afs.FS
	path *filepath.Processor
}

func (s *service) Initialize(req *Request) (*Response, error) {

	// look for relative paths
	template, err := s.path.Abs(req.Template)
	if err != nil {
		return nil, err
	}

	// is this a file or directory?
	stat, err := s.fs.Stat(template)
	if err == nil && stat.IsDir() {
		template = s.path.Join(template, ".caster.yml")
	} else if !errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}

	content := `files:
- name: hello.txt
  content: "hello world"`

	err = s.fs.WriteFile(template, []byte(content), 0666)
	return &Response{}, err
}
