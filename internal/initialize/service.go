package initialize

import (
	"github.com/patrickhuber/go-xplat/filepath"
	"github.com/patrickhuber/go-xplat/fs"
)

type Request struct {
	Template string `yaml:"omitempty"`
}

type Response struct{}

type Service interface {
	Initialize(req *Request) (*Response, error)
}

func NewService(fs fs.FS, path *filepath.Processor) Service {
	return &service{
		fs:   fs,
		path: path,
	}
}

type service struct {
	fs   fs.FS
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
	if err != nil {
		return nil, err
	}
	if stat.IsDir() {
		template = s.path.Join(template, ".caster.yml")
	}

	content := `files:
- name: hello.txt
  content: "hello world"`

	err = s.fs.WriteFile(template, []byte(content), 0666)
	return &Response{}, err
}
