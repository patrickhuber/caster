package initialize_test

import (
	"testing"

	"github.com/patrickhuber/caster/internal/initialize"
	"github.com/patrickhuber/go-xplat/filepath"
	"github.com/patrickhuber/go-xplat/fs"
	"github.com/patrickhuber/go-xplat/os"
	"github.com/patrickhuber/go-xplat/platform"
	"github.com/stretchr/testify/require"
)

func TestService(t *testing.T) {
	o := os.NewMock(os.WithPlatform(platform.Linux))
	path := filepath.NewProcessorWithOS(o)
	fs := fs.NewMemory(fs.WithProcessor(path))

	svc := initialize.NewService(fs, path)

	wd, err := o.WorkingDirectory()
	require.NoError(t, err)
	err = fs.MkdirAll(wd, 0666)
	require.NoError(t, err)

	res, err := svc.Initialize(&initialize.Request{})
	require.NoError(t, err)
	require.NotNil(t, res)

	filePath := path.Join(wd, ".caster.yml")
	ok, err := fs.Exists(filePath)
	require.NoError(t, err)
	require.True(t, ok)

	content, err := fs.ReadFile(filePath)
	require.NoError(t, err)
	require.NotNil(t, content)
}
