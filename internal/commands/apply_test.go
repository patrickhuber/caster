package commands_test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestApply(t *testing.T) {

	t.Run("basic", func(t *testing.T) {
		cx := SetupTestContext(t)
		err := cx.fs.WriteFile("/template/.caster.yml", []byte("files:\n- name: test.txt\n"), 0600)
		require.NoError(t, err)

		args := []string{"caster", "apply", "-t", "/template"}

		err = cx.app.Run(args)
		require.NoError(t, err)

		ok, err := cx.fs.Exists("/working/test.txt")
		require.NoError(t, err)
		require.True(t, ok)
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

		args := []string{"caster", "apply", "-t", "/template"}

		err = cx.app.Run(args)
		require.NoError(t, err)

		ok, err := cx.fs.Exists("/working/test.txt")
		require.NoError(t, err)
		require.True(t, ok)

		content, err := cx.fs.ReadFile("/working/test.txt")
		require.NoError(t, err)
		require.Equal(t, []byte("value"), content)
	})
	t.Run("multi_data", func(t *testing.T) {
		cx := SetupTestContext(t)
		cx.fs.WriteFile("/template/.caster.yml", []byte("files:\n- name: test.txt\n  content: {{.first}}{{.second}}"), 0600)
		cx.fs.WriteFile("/data/1.yml", []byte("first: first"), 0600)
		cx.fs.WriteFile("/data/2.yml", []byte("second: second"), 0600)

		args := []string{"caster", "apply", "--var-file", "/data/1.yml", "--var-file", "/data/2.yml", "-t", "/template"}

		err := cx.app.Run(args)
		require.NoError(t, err)

		ok, err := cx.fs.Exists("/working/test.txt")
		require.NoError(t, err)
		require.True(t, ok)

		want := `firstsecond`
		content, err := cx.fs.ReadFile("/working/test.txt")
		require.NoError(t, err)
		require.Equal(t, []byte(want), content)
	})

	t.Run("multi_arg", func(t *testing.T) {
		cx := SetupTestContext(t)
		cx.fs.WriteFile("/template/.caster.yml", []byte("files:\n- name: test.txt\n  content: {{.first}}{{.second}}"), 0600)

		args := []string{"caster", "apply", "--var", "first=first", "--var", "second=second", "-t", "/template"}

		err := cx.app.Run(args)
		require.NoError(t, err)

		ok, err := cx.fs.Exists("/working/test.txt")
		require.NoError(t, err)
		require.True(t, ok)

		want := `firstsecond`
		content, err := cx.fs.ReadFile("/working/test.txt")
		require.NoError(t, err)
		require.Equal(t, []byte(want), content)
	})
	t.Run("mixed_arg", func(t *testing.T) {
		cx := SetupTestContext(t)
		cx.fs.WriteFile("/template/.caster.yml", []byte("files:\n- name: test.txt\n  content: {{.key}}"), 0600)
		cx.fs.WriteFile("/data/1.yml", []byte("key: first"), 0600)

		args := []string{"caster", "apply", "--var-file", "/data/1.yml", "--var", "key=second", "-t", "/template"}

		err := cx.app.Run(args)
		require.NoError(t, err)

		ok, err := cx.fs.Exists("/working/test.txt")
		require.NoError(t, err)
		require.True(t, ok)

		want := `second`
		content, err := cx.fs.ReadFile("/working/test.txt")
		require.NoError(t, err)
		require.Equal(t, []byte(want), content)
	})
	t.Run("override", func(t *testing.T) {

		cx := SetupTestContext(t)
		cx.fs.WriteFile("/template/.caster.yml", []byte("files:\n- name: test.txt\n  content: {{.key}}"), 0600)
		cx.fs.WriteFile("/data/1.yml", []byte("key: second"), 0600)

		args := []string{"caster", "apply", "--var", "key=first", "--var-file", "/data/1.yml", "-t", "/template"}

		err := cx.app.Run(args)
		require.NoError(t, err)

		ok, err := cx.fs.Exists("/working/test.txt")
		require.NoError(t, err)
		require.True(t, ok)

		want := `second`
		content, err := cx.fs.ReadFile("/working/test.txt")
		require.NoError(t, err)
		require.Equal(t, []byte(want), content)
	})

	t.Run("default", func(t *testing.T) {
		cx := SetupTestContext(t)
		cx.fs.WriteFile("/working/.caster.yml", []byte("files:\n- name: test.txt\n  content: {{.key}}"), 0600)
		cx.fs.WriteFile("/data/1.yml", []byte("key: second"), 0600)

		args := []string{"caster", "apply", "--var", "key=first", "--var-file", "/data/1.yml"}

		err := cx.app.Run(args)
		require.NoError(t, err)

		ok, err := cx.fs.Exists("/working/test.txt")
		require.NoError(t, err)
		require.True(t, ok)

		want := `second`
		content, err := cx.fs.ReadFile("/working/test.txt")
		require.NoError(t, err)
		require.Equal(t, []byte(want), content)
	})
}
