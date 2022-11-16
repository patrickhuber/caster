package interpolate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/patrickhuber/caster/pkg/abstract/env"
	afs "github.com/patrickhuber/caster/pkg/abstract/fs"
	"github.com/patrickhuber/caster/pkg/models"
	"gopkg.in/yaml.v3"
)

type Service interface {
	Interpolate(req *Request) (*Response, error)
}

// NewService creates a new instance of the cast service
func NewService(fs afs.FS, env env.Env) Service {
	return &service{
		fs:  fs,
		env: env,
	}
}

type service struct {
	fs  afs.FS
	env env.Env
}

func (s *service) Interpolate(req *Request) (*Response, error) {

	fileIsSpecified := len(strings.TrimSpace(req.File)) != 0
	directoryIsSpecified := len(strings.TrimSpace(req.Directory)) != 0

	if !fileIsSpecified && !directoryIsSpecified {
		return nil, fmt.Errorf("either source file or source directory must be specified")
	}

	if fileIsSpecified && directoryIsSpecified {
		return nil, fmt.Errorf("source file and source directory are mutually exclusive. Specify one but not both")
	}

	path, err := s.getCasterFile(req)
	if err != nil {
		return nil, err
	}

	content, err := s.getCasterFileContent(path)
	if err != nil {
		return nil, err
	}

	dataMap, err := s.createDataMap(req.Variables)
	if err != nil {
		return nil, err
	}

	rendered, err := s.renderCasterFile(content, path, dataMap)
	if err != nil {
		return nil, err
	}

	structured, err := s.deserializeCasterFile(rendered, filepath.Ext(path))
	if err != nil {
		return nil, err
	}

	return &Response{
		Caster:     *structured,
		SourceFile: path,
	}, nil
}

// createDataMap transforms the variable array to a map[string]any.
// variables are applied in the specified order:
// - files
// - command line arguments
// - environment variables
func (s *service) createDataMap(variables []models.Variable) (map[string]any, error) {
	files := map[string]any{}
	args := map[string]any{}
	env := map[string]any{}
	for _, variable := range variables {
		isEnvVar := len(strings.TrimSpace(variable.Env)) > 0
		isArg := len(strings.TrimSpace(variable.Key)) > 0
		isFile := len(strings.TrimSpace(variable.File)) > 0
		if isArg {
			args[variable.Key] = variable.Value
		}
		if isEnvVar {
			key := strings.TrimPrefix(variable.Env, "CASTER_VAR_")
			env[key] = s.env.Get(variable.Env)
		}
		if isFile {
			content, err := s.fs.Read(variable.File)
			if err != nil {
				return nil, err
			}

			file := map[string]any{}
			err = yaml.Unmarshal(content, file)
			if err != nil {
				return nil, err
			}
			for k, v := range file {
				files[k] = v
			}
		}
	}
	data := files
	for k, v := range args {
		data[k] = v
	}
	for k, v := range env {
		data[k] = v
	}
	return data, nil
}

func (s *service) getCasterFile(req *Request) (string, error) {
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

func (s *service) renderCasterFile(content, sourceFile string, data map[string]interface{}) ([]byte, error) {

	// inject the standard functions defined in sprig
	funcMap := sprig.TxtFuncMap()

	// templatefile renders a template file and then writes the rendered string to the calling template
	funcMap["templatefile"] = func(path string, data interface{}) (string, error) {
		directory := s.fs.Dir(sourceFile)
		content, err := s.fs.Read(s.fs.Join(directory, path))
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

func (s *service) deserializeCasterFile(rendered []byte, extension string) (*models.Caster, error) {
	switch extension {
	case ".yml":
		return s.deserializeYamlCasterFile(rendered)
	case ".json":
		return s.deserializeJsonCasterFile(rendered)
	}
	return nil, fmt.Errorf("unrecognized extension '%s'", extension)
}

func (s *service) deserializeYamlCasterFile(rendered []byte) (*models.Caster, error) {
	var caster models.Caster
	err := yaml.Unmarshal(rendered, &caster)
	return &caster, err
}

func (s *service) deserializeJsonCasterFile(rendered []byte) (*models.Caster, error) {

	var caster models.Caster
	err := json.Unmarshal(rendered, &caster)
	return &caster, err
}
