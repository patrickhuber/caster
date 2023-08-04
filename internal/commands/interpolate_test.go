package commands_test

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/patrickhuber/caster/internal/global"
	"github.com/stretchr/testify/require"
)

func TestInterpolate(t *testing.T) {

	t.Run("basic", func(t *testing.T) {
		cx := SetupTestContext(t)
		err := cx.fs.WriteFile("/template/.caster.yml", []byte("files:\n- name: test.txt\n"), 0600)
		require.NoError(t, err)

		args := []string{"caster", "interpolate", "-t", "/template"}
		cx.app.Metadata[global.OSArgs] = args

		err = cx.app.Run(args)
		require.NoError(t, err)

		buf, ok := cx.console.Out().(*bytes.Buffer)
		require.True(t, ok)
		require.Equal(t, "files:\n  - name: test.txt\n", buf.String())
	})
	t.Run("env", func(t *testing.T) {
		template := `files:
  - name: test.txt
    content: {{ .key }}
`
		cx := SetupTestContext(t)
		err := cx.fs.WriteFile("/template/.caster.yml", []byte(template), 0600)
		require.NoError(t, err)

		cx.env.Set("CASTER_VAR_key", "value")

		args := []string{"caster", "interpolate", "-t", "/template"}
		cx.app.Metadata[global.OSArgs] = args

		err = cx.app.Run(args)
		require.NoError(t, err)

		buf, ok := cx.console.Out().(*bytes.Buffer)
		require.True(t, ok)
		want := `files:
  - name: test.txt
    content: value
`
		have := buf.String()
		require.Equal(t, want, have, cmp.Diff(have, want))
	})
	t.Run("multi_data", func(t *testing.T) {
		cx := SetupTestContext(t)
		cx.fs.WriteFile("/template/.caster.yml", []byte("files:\n- name: test.txt\n  content: {{.first}}{{.second}}"), 0600)
		cx.fs.WriteFile("/data/1.yml", []byte("first: first"), 0600)
		cx.fs.WriteFile("/data/2.yml", []byte("second: second"), 0600)

		args := []string{"caster", "interpolate", "--var-file", "/data/1.yml", "--var-file", "/data/2.yml", "-t", "/template"}
		cx.app.Metadata[global.OSArgs] = args
		err := cx.app.Run(args)
		require.NoError(t, err)

		buf, ok := cx.console.Out().(*bytes.Buffer)
		require.True(t, ok)
		want := `files:
  - name: test.txt
    content: firstsecond
`
		have := buf.String()
		require.Equal(t, want, have, cmp.Diff(have, want))
	})
	t.Run("multi_arg", func(t *testing.T) {
		cx := SetupTestContext(t)
		cx.fs.WriteFile("/template/.caster.yml", []byte("files:\n- name: test.txt\n  content: {{.first}}{{.second}}"), 0600)

		args := []string{"caster", "interpolate", "--var", "first=first", "--var", "second=second", "-t", "/template"}
		cx.app.Metadata[global.OSArgs] = args
		err := cx.app.Run(args)
		require.NoError(t, err)

		buf, ok := cx.console.Out().(*bytes.Buffer)
		require.True(t, ok)
		want := `files:
  - name: test.txt
    content: firstsecond
`
		have := buf.String()
		require.Equal(t, want, have, cmp.Diff(have, want))
	})
	t.Run("mixed_arg", func(t *testing.T) {
		cx := SetupTestContext(t)
		cx.fs.WriteFile("/template/.caster.yml", []byte("files:\n- name: test.txt\n  content: {{.key}}"), 0600)
		cx.fs.WriteFile("/data/1.yml", []byte("key: first"), 0600)

		args := []string{"caster", "interpolate", "--var-file", "/data/1.yml", "--var", "key=second", "-t", "/template"}
		cx.app.Metadata[global.OSArgs] = args
		err := cx.app.Run(args)
		require.NoError(t, err)

		buf, ok := cx.console.Out().(*bytes.Buffer)
		require.True(t, ok)
		want := `files:
  - name: test.txt
    content: second
`
		have := buf.String()
		require.Equal(t, want, have, cmp.Diff(have, want))
	})
	t.Run("override", func(t *testing.T) {

		cx := SetupTestContext(t)
		cx.fs.WriteFile("/template/.caster.yml", []byte("files:\n- name: test.txt\n  content: {{.key}}"), 0600)
		cx.fs.WriteFile("/data/1.yml", []byte("key: second"), 0600)

		args := []string{"caster", "interpolate", "--var", "key=first", "--var-file", "/data/1.yml", "-t", "/template"}
		cx.app.Metadata[global.OSArgs] = args
		err := cx.app.Run(args)
		require.NoError(t, err)

		buf, ok := cx.console.Out().(*bytes.Buffer)
		require.True(t, ok)
		want := `files:
  - name: test.txt
    content: second
`
		have := buf.String()
		require.Equal(t, want, have, cmp.Diff(have, want))
	})

	t.Run("default", func(t *testing.T) {
		cx := SetupTestContext(t)
		cx.fs.WriteFile("/working/.caster.yml", []byte("files:\n- name: test.txt\n  content: {{.key}}"), 0600)
		cx.fs.WriteFile("/data/1.yml", []byte("key: second"), 0600)

		args := []string{"caster", "interpolate", "--var", "key=first", "--var-file", "/data/1.yml"}
		cx.app.Metadata[global.OSArgs] = args
		err := cx.app.Run(args)
		require.NoError(t, err)

		buf, ok := cx.console.Out().(*bytes.Buffer)
		require.True(t, ok)
		want := `files:
  - name: test.txt
    content: second
`
		have := buf.String()
		require.Equal(t, want, have, cmp.Diff(have, want))
	})
}
