package cast_test

import (
	"strings"
	"testing"

	"github.com/patrickhuber/caster/internal/cast"
	"github.com/patrickhuber/caster/internal/interpolate"
	"github.com/patrickhuber/caster/internal/models"
	"github.com/stretchr/testify/require"

	"github.com/patrickhuber/go-xplat/arch"
	"github.com/patrickhuber/go-xplat/host"
	"github.com/patrickhuber/go-xplat/platform"
)

func Setup(t *testing.T, h *host.Host, content []byte, inter interpolate.Service, request *cast.Request) {
	// a template is either a path to a file or directory
	// we need to take the path and determine its type and generate the appropriate
	// test file system
	template := request.Template

	// if the template is completely empty, use a default path
	if len(strings.TrimSpace(template)) == 0 {
		template = "/template"
	}

	// a file will have an extension
	isFile := len(strings.TrimSpace(h.Path.Ext(template))) > 0
	if !isFile {
		// this is a directory so we need to append the default file to the directory
		template = h.Path.Join(template, ".caster.yml")
	}

	err := h.FS.MkdirAll("/output", 0600)
	require.NoError(t, err)

	err = h.FS.MkdirAll("/template", 0600)
	require.NoError(t, err)

	err = h.FS.WriteFile(template, content, 0600)
	require.NoError(t, err)

	source := h.Path.Dir(template)
	sourceInfo, err := h.FS.Stat(source)
	require.NoError(t, err)
	require.True(t, sourceInfo.IsDir())

	svc := cast.NewService(h.FS, inter, h.Path)

	err = svc.Cast(request)
	require.NoError(t, err)

	AssertExists(t, h, request.Target)
}

func AssertExists(t *testing.T, h *host.Host, path string) {
	ok, err := h.FS.Exists(path)
	require.NoError(t, err)
	require.True(t, ok, "expected '%s' to exist", path)
}

func AssertContents(t *testing.T, h *host.Host, path, content string) {
	data, err := h.FS.ReadFile(path)
	require.NoError(t, err)
	require.Equal(t, content, string(data))
}

func TestService(t *testing.T) {
	type file struct {
		path    string
		content string
		dir     bool
	}
	type test struct {
		name     string
		template string
		files    []file
		request  *cast.Request
		stage    []file
	}
	tests := []test{
		{
			"apply_file", `files:
- name: test.yml
  content: "test: test"
folders:
- name: sub
  files: 
  - name: test.yml`,
			[]file{
				{"/output", "", true},
				{"/output/test.yml", "test: test", false},
				{"/output", "", true},
				{"/output/sub/test.yml", "", false},
			}, &cast.Request{
				Template: "/template/custom.yml",
				Target:   "/output",
			}, []file{},
		},
		{
			"replaces_file_names",
			`---
files:
- name: {{"hello"}}{{"world"}}.yml
  content: "hello: world"`,
			[]file{
				{"/output", "", true},
				{"/output/helloworld.yml", "hello: world", false},
			},
			&cast.Request{Template: "/template", Target: "/output"},
			[]file{},
		},
		{
			"replaces_folder_names",
			`---
folders:
- name: {{"hello"}}
  files:
  - name: 1.yml
    content: "one: 1"`,
			[]file{
				{"/output/hello", "", true},
				{"/output/hello/1.yml", "one: 1", false},
			},
			&cast.Request{Template: "/template", Target: "/output"},
			[]file{},
		},
		{
			"ref",
			`---
files:
- name: test
  ref: test.txt`,
			[]file{
				{"/output/hello", "", true},
				{"/output/hello/1.yml", "one: 1", false},
			},
			&cast.Request{Template: "/template", Target: "/output"},
			[]file{{"/template/test.txt", "test", false}},
		},
		{
			"content",
			`---
files:
- name: test
  content: {{ templatefile "./test.txt" . }}`,
			[]file{
				{"/output/test", "value", false},
			},
			&cast.Request{
				Template: "/template",
				Target:   "/output",
				Variables: []models.Variable{
					{Key: "key", Value: "value"},
				},
			},
			[]file{{"/template/test.txt", "{{ .key }}", false}},
		},
		{
			"multi",
			`files:
- name: test
  content: |
    {{- templatefile "./test.txt" . | nindent 4 }}`,
			[]file{
				{"/output/test", "value\nvalue", false},
			},
			&cast.Request{
				Template: "/template",
				Target:   "/output",
				Variables: []models.Variable{
					{Key: "key", Value: "value"},
				},
			},
			[]file{{"/template/test.txt", "{{ .key }}", false}},
		},
		{
			"varfile",
			`---
files:
- name: test.yml
  content: {{ .variable }}`,
			[]file{
				{"/output/test.yml", "test", false},
			},
			&cast.Request{
				Template: "/template",
				Target:   "/output",
				Variables: []models.Variable{
					{Key: "key", Value: "value"},
				},
			},
			[]file{{"/data.yml", "variable: test", false}},
		},
		{
			"varfile",
			`---
files:
- name: test.yml
  content: {{ .variable }}`,
			[]file{
				{"/output/test.yml", "test", true},
			},
			&cast.Request{
				Template: "/template",
				Target:   "/output",
				Variables: []models.Variable{
					{Key: "variable", Value: "test"},
				},
			},
			[]file{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := host.NewTest(platform.Linux, arch.AMD64)
			h.OS.ChangeDirectory("/")
			svc := interpolate.NewService(h.FS, h.Env, h.Path)
			for _, file := range test.stage {
				err := h.FS.WriteFile(file.path, []byte(file.content), 0666)
				require.NoError(t, err)
			}
			Setup(t, h, []byte(test.template), svc, test.request)
			for _, file := range test.files {
				AssertExists(t, h, file.path)
				if !file.dir {
					AssertContents(t, h, file.path, file.content)
				}
			}
		})
	}
}
