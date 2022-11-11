package cast

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	afs "github.com/patrickhuber/caster/pkg/abstract/fs"
	"gopkg.in/yaml.v3"
)

// Service handles casting of a template
type Service interface {
	Cast(req *CastRequest) error
}

type service struct {
	fs afs.FS
}

// NewService creates a new instance of the cast service
func NewService(fs afs.FS) Service {
	return &service{
		fs: fs,
	}
}

func (s *service) Cast(req *CastRequest) error {
	fileIsSpecified := len(strings.TrimSpace(req.File)) != 0
	directoryIsSpecified := len(strings.TrimSpace(req.Directory)) != 0
	targetIsSpecified := len(strings.TrimSpace(req.Target)) != 0

	if !fileIsSpecified && !directoryIsSpecified {
		return fmt.Errorf("either source file or source directory must be specified")
	}

	if fileIsSpecified && directoryIsSpecified {
		return fmt.Errorf("source file and source directory are mutually exclusive. Specify one but not both")
	}

	if !targetIsSpecified {
		return fmt.Errorf("target must be specified")
	}

	path, err := s.getCasterFile(req)
	if err != nil {
		return err
	}

	source := s.fs.Dir(path)

	content, err := s.getCasterFileContent(path)
	if err != nil {
		return err
	}

	rendered, err := s.renderCasterFile(content, req.Data)
	if err != nil {
		return err
	}

	structured, err := s.deserializeCasterFile(rendered, filepath.Ext(path))
	if err != nil {
		return err
	}

	return s.executeCasterFile(structured, source, req.Target)
}

func (s *service) getCasterFile(req *CastRequest) (string, error) {
	file := strings.TrimSpace(req.File)
	if len(file) > 0 {
		return file, nil
	}

	// read the caster file in the directory
	// if it doesn't exist, return an error saying not found
	files, err := s.fs.ReadDirRegex(req.Directory, "[.]caster[.](yml|json)")
	if err != nil {
		return "", err
	}
	if len(files) == 0 {
		return "", fmt.Errorf("template folder '%s' missing .caster.(yml|json) file", req.Directory)
	}

	file = s.fs.Join(req.Directory, files[0].Name())
	return file, nil
}

func (s *service) getCasterFileContent(path string) (string, error) {
	content, err := s.fs.Read(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func (s *service) renderCasterFile(content string, data map[string]interface{}) ([]byte, error) {

	// inject the standard functions defined in sprig
	funcMap := sprig.TxtFuncMap()

	// templatefile renders a template file and then writes the rendered string to the calling template
	funcMap["templatefile"] = func(path string, data interface{}) (string, error) {
		content, err := s.fs.Read(path)
		if err != nil {
			return "", err
		}
		t, err := template.
			New("inner").
			Funcs(sprig.TxtFuncMap()).
			Parse(string(content))
		if err != nil {
			return "", err
		}
		var writer bytes.Buffer
		err = t.Execute(&writer, data)
		return writer.String(), err
	}

	// parse the template
	t, err := template.
		New("caster").
		Funcs(funcMap).
		Parse(content)
	if err != nil {
		return nil, err
	}

	// execute the template
	var writer bytes.Buffer
	err = t.Execute(&writer, data)
	return writer.Bytes(), err
}

func (s *service) deserializeCasterFile(rendered []byte, extension string) (*Caster, error) {
	switch extension {
	case ".yml":
		return s.deserializeYamlCasterFile(rendered)
	case ".json":
		return s.deserializeJsonCasterFile(rendered)
	}
	return nil, fmt.Errorf("unrecognized extension '%s'", extension)
}

func (s *service) deserializeYamlCasterFile(rendered []byte) (*Caster, error) {
	var caster Caster
	err := yaml.Unmarshal(rendered, &caster)
	return &caster, err
}

func (s *service) deserializeJsonCasterFile(rendered []byte) (*Caster, error) {

	var caster Caster
	err := json.Unmarshal(rendered, &caster)
	return &caster, err
}

func (s *service) executeCasterFile(caster *Caster, source, target string) error {
	err := s.castFiles(source, target, source, caster.Files)
	if err != nil {
		return err
	}
	ok, err := s.fs.Exists(target)
	if err != nil {
		return err
	}
	if !ok {
		err := s.fs.Mkdir(target, 0600)
		if err != nil {
			return err
		}
	}
	return s.castFolders(source, target, source, caster.Folders)
}

func (s *service) castFolders(source, target, path string, folders []Folder) error {
	for _, folder := range folders {
		sourcePath := s.fs.Join(path, folder.Name)
		err := s.castFolder(&folder, source, target, sourcePath)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *service) castFolder(folder *Folder, source, target, path string) error {
	rel, err := s.fs.Rel(source, path)
	if err != nil {
		return err
	}

	targetPath := s.fs.Join(target, rel)
	err = s.fs.Mkdir(targetPath, 0600)
	if err != nil {
		return err
	}

	err = s.castFiles(source, target, path, folder.Files)
	if err != nil {
		return err
	}

	return s.castFolders(source, target, path, folder.Folders)
}

func (s *service) castFiles(source, target string, path string, files []File) error {
	for _, file := range files {
		sourcePath := s.fs.Join(path, file.Name)
		err := s.castFile(&file, source, target, sourcePath)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *service) castFile(file *File, source, target, path string) error {
	rel, err := s.fs.Rel(source, path)
	if err != nil {
		return err
	}
	targetPath := s.fs.Join(target, rel)
	content := []byte(file.Content)

	// is the ref set and the content empty?
	if file.Content == "" && file.Ref != "" {
		path := s.fs.Join(source, file.Ref)
		content, err = s.fs.Read(path)
		if err != nil {
			return err
		}
	}
	return s.fs.Write(targetPath, content, 0600)
}
