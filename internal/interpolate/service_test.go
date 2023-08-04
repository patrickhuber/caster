package interpolate_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/patrickhuber/caster/internal/interpolate"
	"github.com/patrickhuber/caster/internal/models"
	"github.com/patrickhuber/go-xplat/env"
	"github.com/patrickhuber/go-xplat/filepath"
	afs "github.com/patrickhuber/go-xplat/fs"
	"github.com/patrickhuber/go-xplat/os"
	"github.com/patrickhuber/go-xplat/platform"
)

type ServiceTestContext struct {
	fs   afs.FS
	e    env.Environment
	svc  interpolate.Service
	path *filepath.Processor
}

func TestService(t *testing.T) {
	t.Run("can interpolate", func(t *testing.T) {
		cx := CreateServiceTestContext(t)
		template := `---
files:
- name: test.txt
  content: test
`
		err := cx.fs.WriteFile("/template/.caster.yml", []byte(template), 0600)
		require.NoError(t, err)

		resp, err := cx.svc.Interpolate(&interpolate.Request{
			Template: "/template",
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotEqual(t, models.Caster{}, resp.Caster)
		require.Equal(t, 1, len(resp.Caster.Files))

		file := resp.Caster.Files[0]
		require.Equal(t, "test", file.Content)
		require.Equal(t, "test.txt", file.Name)
	})
	t.Run("can interpolate with data", func(t *testing.T) {
		cx := CreateServiceTestContext(t)
		template := `---
files:
- name: test.txt
  content: {{ .key }}
`
		err := cx.fs.WriteFile("/template/.caster.yml", []byte(template), 0600)
		require.NoError(t, err)

		resp, err := cx.svc.Interpolate(&interpolate.Request{
			Template: "/template",
			Variables: []models.Variable{
				{
					Key:   "key",
					Value: "value",
				},
			},
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotEqual(t, models.Caster{}, resp.Caster)
		require.Equal(t, 1, len(resp.Caster.Files))

		file := resp.Caster.Files[0]
		require.Equal(t, "value", file.Content)
		require.Equal(t, "test.txt", file.Name)
	})
}

func CreateServiceTestContext(t *testing.T) *ServiceTestContext {
	o := os.NewMock(os.WithPlatform(platform.Linux))
	path := filepath.NewProcessorWithOS(o)
	fs := afs.NewMemory(afs.WithProcessor(path))
	require.NoError(t, fs.Mkdir("/", 0600))
	require.NoError(t, fs.Mkdir("/template", 0600))
	e := env.NewMemory()
	svc := interpolate.NewService(fs, e, path)
	return &ServiceTestContext{
		fs:   fs,
		path: path,
		e:    e,
		svc:  svc,
	}
}
