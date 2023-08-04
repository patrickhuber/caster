package commands_test

import (
	"testing"

	"github.com/patrickhuber/caster/internal/global"
	"github.com/stretchr/testify/require"
)

func TestInitialize(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		cx := SetupTestContext(t)

		args := []string{"caster", "init"}
		cx.app.Metadata[global.OSArgs] = args

		err := cx.app.Run(args)
		require.NoError(t, err)

		wd, err := cx.os.WorkingDirectory()
		require.NoError(t, err)

		ok, err := cx.fs.Exists(cx.path.Join(wd, ".caster.yml"))
		require.NoError(t, err)
		require.True(t, ok)
	})

	t.Run("template_dir", func(t *testing.T) {
		cx := SetupTestContext(t)

		args := []string{"caster", "init", "-t", "/template"}
		cx.app.Metadata[global.OSArgs] = args

		err := cx.app.Run(args)
		require.NoError(t, err)

		ok, err := cx.fs.Exists(cx.path.Join("/template", ".caster.yml"))
		require.NoError(t, err)
		require.True(t, ok)
	})

	t.Run("template_file", func(t *testing.T) {
		cx := SetupTestContext(t)

		args := []string{"caster", "init", "-t", "/template/test.yml"}
		cx.app.Metadata[global.OSArgs] = args

		err := cx.app.Run(args)
		require.NoError(t, err)

		ok, err := cx.fs.Exists(cx.path.Join("/template", "test.yml"))
		require.NoError(t, err)
		require.True(t, ok)
	})
}
