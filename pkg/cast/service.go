package cast

import (
	"fmt"
	"strings"

	afs "github.com/patrickhuber/caster/pkg/abstract/fs"
	"github.com/patrickhuber/caster/pkg/interpolate"
	"github.com/patrickhuber/caster/pkg/models"
)

// Service handles casting of a template
type Service interface {
	Cast(req *Request) error
}

type service struct {
	fs    afs.FS
	inter interpolate.Service
}

// NewService creates a new instance of the cast service
func NewService(fs afs.FS, inter interpolate.Service) Service {
	return &service{
		fs:    fs,
		inter: inter,
	}
}

func (s *service) Cast(req *Request) error {

	variables := []models.Variable{}
	for _, v := range req.Variables {
		variables = append(variables, models.Variable{
			Env:   v.Env,
			File:  v.File,
			Key:   v.Key,
			Value: v.Value,
		})
	}
	resp, err := s.inter.Interpolate(&interpolate.Request{
		File:      req.File,
		Directory: req.Directory,
		Variables: variables,
	})

	if err != nil {
		return err
	}

	targetIsSpecified := len(strings.TrimSpace(req.Target)) != 0

	if !targetIsSpecified {
		return fmt.Errorf("target must be specified")
	}

	source := s.fs.Dir(resp.SourceFile)
	return s.executeCasterFile(&resp.Caster, source, req.Target)
}

func (s *service) executeCasterFile(caster *models.Caster, source, target string) error {
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

func (s *service) castFolders(source, target, path string, folders []models.Folder) error {
	for _, folder := range folders {
		sourcePath := s.fs.Join(path, folder.Name)
		err := s.castFolder(&folder, source, target, sourcePath)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *service) castFolder(folder *models.Folder, source, target, path string) error {
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

func (s *service) castFiles(source, target string, path string, files []models.File) error {
	for _, file := range files {
		sourcePath := s.fs.Join(path, file.Name)
		err := s.castFile(&file, source, target, sourcePath)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *service) castFile(file *models.File, source, target, path string) error {
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
