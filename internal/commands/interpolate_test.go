package commands_test

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/patrickhuber/caster/internal/commands"
	"github.com/patrickhuber/caster/internal/global"
	"github.com/patrickhuber/caster/internal/setup"
	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/go-xplat/console"
	"github.com/patrickhuber/go-xplat/env"
	"github.com/patrickhuber/go-xplat/fs"
	"github.com/patrickhuber/go-xplat/os"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

type InterpolateTestContext struct {
	app       *cli.App
	container di.Container
	con       console.Console
	f         fs.FS
	e         env.Environment
}

func TestInterpolate(t *testing.T) {

	t.Run("basic", func(t *testing.T) {
		cx := SetupInterpolateTestContext(t)
		err := cx.f.WriteFile("/template/.caster.yml", []byte("files:\n- name: test.txt\n"), 0600)
		require.NoError(t, err)

		args := []string{"caster", "interpolate", "-t", "/template"}
		cx.app.Metadata[global.OSArgs] = args

		err = cx.app.Run(args)
		require.NoError(t, err)

		buf, ok := cx.con.Out().(*bytes.Buffer)
		require.True(t, ok)
		require.Equal(t, "files:\n  - name: test.txt\n", buf.String())
	})
	t.Run("env", func(t *testing.T) {
		template := `files:
  - name: test.txt
    content: {{ .key }}
`
		cx := SetupInterpolateTestContext(t)
		err := cx.f.WriteFile("/template/.caster.yml", []byte(template), 0600)
		require.NoError(t, err)

		cx.e.Set("CASTER_VAR_key", "value")

		args := []string{"caster", "interpolate", "-t", "/template"}
		cx.app.Metadata[global.OSArgs] = args

		err = cx.app.Run(args)
		require.NoError(t, err)

		buf, ok := cx.con.Out().(*bytes.Buffer)
		require.True(t, ok)
		want := `files:
  - name: test.txt
    content: value
`
		have := buf.String()
		require.Equal(t, want, have, cmp.Diff(have, want))
	})
	t.Run("multi_data", func(t *testing.T) {
		cx := SetupInterpolateTestContext(t)
		cx.f.WriteFile("/template/.caster.yml", []byte("files:\n- name: test.txt\n  content: {{.first}}{{.second}}"), 0600)
		cx.f.WriteFile("/data/1.yml", []byte("first: first"), 0600)
		cx.f.WriteFile("/data/2.yml", []byte("second: second"), 0600)

		args := []string{"caster", "interpolate", "--var-file", "/data/1.yml", "--var-file", "/data/2.yml", "-t", "/template"}
		cx.app.Metadata[global.OSArgs] = args
		err := cx.app.Run(args)
		require.NoError(t, err)

		buf, ok := cx.con.Out().(*bytes.Buffer)
		require.True(t, ok)
		want := `files:
  - name: test.txt
    content: firstsecond
`
		have := buf.String()
		require.Equal(t, want, have, cmp.Diff(have, want))
	})
	t.Run("multi_arg", func(t *testing.T) {
		cx := SetupInterpolateTestContext(t)
		cx.f.WriteFile("/template/.caster.yml", []byte("files:\n- name: test.txt\n  content: {{.first}}{{.second}}"), 0600)

		args := []string{"caster", "interpolate", "--var", "first=first", "--var", "second=second", "-t", "/template"}
		cx.app.Metadata[global.OSArgs] = args
		err := cx.app.Run(args)
		require.NoError(t, err)

		buf, ok := cx.con.Out().(*bytes.Buffer)
		require.True(t, ok)
		want := `files:
  - name: test.txt
    content: firstsecond
`
		have := buf.String()
		require.Equal(t, want, have, cmp.Diff(have, want))
	})
	t.Run("mixed_arg", func(t *testing.T) {
		cx := SetupInterpolateTestContext(t)
		cx.f.WriteFile("/template/.caster.yml", []byte("files:\n- name: test.txt\n  content: {{.key}}"), 0600)
		cx.f.WriteFile("/data/1.yml", []byte("key: first"), 0600)

		args := []string{"caster", "interpolate", "--var-file", "/data/1.yml", "--var", "key=second", "-t", "/template"}
		cx.app.Metadata[global.OSArgs] = args
		err := cx.app.Run(args)
		require.NoError(t, err)

		buf, ok := cx.con.Out().(*bytes.Buffer)
		require.True(t, ok)
		want := `files:
  - name: test.txt
    content: second
`
		have := buf.String()
		require.Equal(t, want, have, cmp.Diff(have, want))
	})
	t.Run("override", func(t *testing.T) {

		cx := SetupInterpolateTestContext(t)
		cx.f.WriteFile("/template/.caster.yml", []byte("files:\n- name: test.txt\n  content: {{.key}}"), 0600)
		cx.f.WriteFile("/data/1.yml", []byte("key: second"), 0600)

		args := []string{"caster", "interpolate", "--var", "key=first", "--var-file", "/data/1.yml", "-t", "/template"}
		cx.app.Metadata[global.OSArgs] = args
		err := cx.app.Run(args)
		require.NoError(t, err)

		buf, ok := cx.con.Out().(*bytes.Buffer)
		require.True(t, ok)
		want := `files:
  - name: test.txt
    content: second
`
		have := buf.String()
		require.Equal(t, want, have, cmp.Diff(have, want))
	})

	t.Run("default", func(t *testing.T) {
		cx := SetupInterpolateTestContext(t)
		cx.f.WriteFile("/working/.caster.yml", []byte("files:\n- name: test.txt\n  content: {{.key}}"), 0600)
		cx.f.WriteFile("/data/1.yml", []byte("key: second"), 0600)

		args := []string{"caster", "interpolate", "--var", "key=first", "--var-file", "/data/1.yml"}
		cx.app.Metadata[global.OSArgs] = args
		err := cx.app.Run(args)
		require.NoError(t, err)

		buf, ok := cx.con.Out().(*bytes.Buffer)
		require.True(t, ok)
		want := `files:
  - name: test.txt
    content: second
`
		have := buf.String()
		require.Equal(t, want, have, cmp.Diff(have, want))
	})
}

func SetupInterpolateTestContext(t *testing.T) *InterpolateTestContext {
	var err error
	s := setup.NewTest()
	container := s.Container()

	o, err := di.Resolve[os.OS](container)
	require.NoError(t, err)

	wd, err := o.WorkingDirectory()
	require.NoError(t, err)

	con, err := di.Resolve[console.Console](container)
	require.NoError(t, err)

	f, err := di.Resolve[fs.FS](container)
	require.NoError(t, err)

	paths := []string{"/", "/template", "/data", wd, o.Home()}
	for _, p := range paths {
		err = f.MkdirAll(p, 0666)
		require.NoError(t, err)
	}

	e, err := di.Resolve[env.Environment](container)
	require.NoError(t, err)

	app := &cli.App{
		Version: "1.0.0",
		Metadata: map[string]interface{}{
			global.DependencyInjectionContainer: container,
		},
		Commands: []*cli.Command{
			commands.Interpolate,
		},
		Reader:    con.In(),
		ErrWriter: con.Error(),
		Writer:    con.Out(),
	}

	return &InterpolateTestContext{
		app:       app,
		container: container,
		con:       con,
		f:         f,
		e:         e,
	}
}
