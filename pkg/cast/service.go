package cast

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	afs "github.com/patrickhuber/caster/pkg/abstract/fs"
	"gopkg.in/yaml.v3"
)

// Service handles casting of a template
type Service interface {
	Cast(source string, target string, data map[string]interface{}) error
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

func (s *service) Cast(source, target string, data map[string]interface{}) error {

	info, err := s.getCasterFile(source)
	if err != nil {
		return err
	}

	content, err := s.getCasterFileContent(source, info)
	if err != nil {
		return err
	}

	rendered, err := s.renderCasterFile(content, data)
	if err != nil {
		return err
	}

	structured, err := s.deserializeCasterFile(rendered, filepath.Ext(info.Name()))
	if err != nil {
		return err
	}

	return s.executeCasterFile(structured, source, target)
}

func (s *service) getCasterFile(source string) (os.FileInfo, error) {
	// read the caster file in the directory
	// if it doesn't exist, return an error saying not found
	files, err := s.fs.ReadDirRegex(source, "[.]caster[.](yml|json)")
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("template folder '%s' missing .caster.(yml|json) file", source)
	}

	return files[0], nil
}

func (s *service) getCasterFileContent(source string, casterFile os.FileInfo) (string, error) {
	content, err := s.fs.Read(s.fs.Join(source, casterFile.Name()))
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
