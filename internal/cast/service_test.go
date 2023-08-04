package cast_test

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/caster/internal/cast"
	"github.com/patrickhuber/caster/internal/interpolate"
	"github.com/patrickhuber/caster/internal/models"
	"gopkg.in/yaml.v3"

	"github.com/patrickhuber/go-xplat/env"
	"github.com/patrickhuber/go-xplat/filepath"
	afs "github.com/patrickhuber/go-xplat/fs"
	"github.com/patrickhuber/go-xplat/os"
)

type ServiceTest interface {
	Setup(template *models.Caster, request *cast.Request)
	SetupString(content string, request *cast.Request)
	SetupBytes(content []byte, request *cast.Request)
	AssertExists(path string)
	AssertContents(path, content string)
	FileSystem() afs.FS
	Environment() env.Environment
}

type serviceTest struct {
	fs    afs.FS
	path  *filepath.Processor
	env   env.Environment
	inter interpolate.Service
}

func NewServiceTest() ServiceTest {
	o := os.NewLinuxMock()
	path := filepath.NewProcessorWithOS(o)
	fs := afs.NewMemory(afs.WithProcessor(path))
	e := env.NewMemory()

	return &serviceTest{
		fs:    fs,
		env:   e,
		inter: interpolate.NewService(fs, e, path),
		path:  path,
	}
}

func (t *serviceTest) Setup(template *models.Caster, request *cast.Request) {
	content, err := yaml.Marshal(template)
	Expect(err).To(BeNil())
	t.SetupBytes(content, request)
}

func (t *serviceTest) SetupString(content string, request *cast.Request) {
	t.SetupBytes([]byte(content), request)
}

func (t *serviceTest) SetupBytes(content []byte, request *cast.Request) {
	// a template is either a path to a file or directory
	// we need to take the path and determine its type and generate the appropriate
	// test file system
	template := request.Template

	// if the template is completely empty, use a default path
	if len(strings.TrimSpace(template)) == 0 {
		template = "/template"
	}

	// a file will have an extension
	isFile := len(strings.TrimSpace(t.path.Ext(template))) > 0
	if !isFile {
		// this is a directory so we need to append the default file to the directory
		template = t.path.Join(template, ".caster.yml")
	}

	err := t.fs.Mkdir("/", 0600)
	Expect(err).To(BeNil())

	err = t.fs.Mkdir("/template", 0600)
	Expect(err).To(BeNil())

	err = t.fs.WriteFile(template, content, 0600)
	Expect(err).To(BeNil())

	source := t.path.Dir(template)
	sourceInfo, err := t.fs.Stat(source)
	Expect(err).To(BeNil())
	Expect(sourceInfo.IsDir()).To(BeTrue())

	svc := cast.NewService(t.fs, t.inter, t.path)

	err = svc.Cast(request)
	Expect(err).To(BeNil())

	t.AssertExists(request.Target)
}

func (t *serviceTest) AssertExists(path string) {
	ok, err := t.fs.Exists(path)
	Expect(err).To(BeNil())
	Expect(ok).To(BeTrue(), "expected '%s' to exist", path)
}

func (t *serviceTest) AssertContents(path, content string) {
	data, err := t.fs.ReadFile(path)
	Expect(err).To(BeNil())
	Expect(string(data)).To(Equal(content))
}

func (t *serviceTest) FileSystem() afs.FS {
	return t.fs
}

func (t *serviceTest) Environment() env.Environment {
	return t.env
}

var _ = Describe("Service", func() {
	Describe("Cast", func() {
		When("caster file specified", func() {
			It("applies from specified file", func() {
				template := &models.Caster{
					Files: []models.File{
						{
							Name:    "test.yml",
							Content: "test: test",
						},
					},
					Folders: []models.Folder{
						{
							Name: "sub",
							Files: []models.File{
								{
									Name: "test.yml",
								},
							},
						},
					},
				}
				t := NewServiceTest()
				t.Setup(template, &cast.Request{
					Template: "/template/custom.yml",
					Target:   "/output",
				})

				t.AssertExists("/output/test.yml")
				t.AssertContents("/output/test.yml", "test: test")
				t.AssertExists("/output/sub")
				t.AssertExists("/output/sub/test.yml")
			})
		})
		It("writes plain files to target", func() {
			template := &models.Caster{
				Files: []models.File{
					{
						Name:    "test.yml",
						Content: "test: test",
					},
				},
				Folders: []models.Folder{
					{
						Name: "sub",
						Files: []models.File{
							{
								Name: "test.yml",
							},
						},
					},
				},
			}
			t := NewServiceTest()
			t.Setup(template, &cast.Request{
				Template: "/template",
				Target:   "/output"})
			t.AssertExists("/output/test.yml")
			t.AssertContents("/output/test.yml", "test: test")
			t.AssertExists("/output/sub")
			t.AssertExists("/output/sub/test.yml")
		})
		It("evaluates file names", func() {
			template := `---
files:
- name: {{"hello"}}{{"world"}}.yml
  content: "hello: world"`

			t := NewServiceTest()
			t.SetupString(template, &cast.Request{Template: "/template", Target: "/output"})
			t.AssertExists("/output")
			t.AssertExists("/output/helloworld.yml")
			t.AssertContents("/output/helloworld.yml", "hello: world")
		})
		It("evaluates folder names", func() {
			template := `---
folders:
- name: {{"hello"}}
  files:
  - name: 1.yml
    content: "one: 1"`

			t := NewServiceTest()
			t.SetupString(template, &cast.Request{Template: "/template", Target: "/output"})
			t.AssertExists("/output/hello")
			t.AssertExists("/output/hello/1.yml")
			t.AssertContents("/output/hello/1.yml", "one: 1")
		})
		It("writes ref", func() {
			template := `---
files:
- name: test
  ref: test.txt`
			t := NewServiceTest()

			err := t.FileSystem().WriteFile("/template/test.txt", []byte("test"), 0644)
			Expect(err).To(BeNil())

			t.SetupString(template, &cast.Request{Template: "/template", Target: "/output"})
			t.AssertExists("/output/test")
			t.AssertContents("/output/test", "test")
		})
		It("interpolates content", func() {
			template := `---
files:
- name: test
  content: {{ templatefile "./test.txt" . }}`

			t := NewServiceTest()
			err := t.FileSystem().WriteFile("/template/test.txt", []byte("{{ .key }}"), 0644)
			Expect(err).To(BeNil())

			t.SetupString(template, &cast.Request{
				Template: "/template",
				Target:   "/output",
				Variables: []models.Variable{
					{Key: "key", Value: "value"},
				},
			})
			t.AssertExists("/output/test")
			t.AssertContents("/output/test", "value")
		})
		It("can indent with multi line string", func() {
			template := `
files:
- name: test
  content: |
    {{- templatefile "./test.txt" . | nindent 4 }}`
			t := NewServiceTest()
			err := t.FileSystem().WriteFile("/template/test.txt", []byte("{{ .key }}\n{{ .key }}"), 0644)
			Expect(err).To(BeNil())

			t.SetupString(template, &cast.Request{
				Template: "/template",
				Target:   "/output",
				Variables: []models.Variable{
					{Key: "key", Value: "value"},
				},
			})
			t.AssertExists("/output/test")
			t.AssertContents("/output/test", "value\nvalue")

		})
		It("can accept variable from file", func() {
			template := `---
files:
- name: test.yml
  content: {{ .variable }}`
			data := "variable: test"

			t := NewServiceTest()
			fs := t.FileSystem()
			err := fs.WriteFile("/data.yml", []byte(data), 0644)
			Expect(err).To(BeNil())

			t.SetupString(template, &cast.Request{
				Template: "/template",
				Target:   "/output",
				Variables: []models.Variable{
					{
						File: "/data.yml",
					},
				}})

			t.AssertExists("/output/test.yml")
			t.AssertContents("/output/test.yml", "test")
		})
		It("can accept variable from arg", func() {
			template := `---
files:
- name: test.yml
  content: {{ .variable }}`
			t := NewServiceTest()
			t.SetupString(template, &cast.Request{
				Template: "/template",
				Target:   "/output",
				Variables: []models.Variable{
					{
						Key:   "variable",
						Value: "test",
					},
				}})

			t.AssertExists("/output/test.yml")
			t.AssertContents("/output/test.yml", "test")
		})
		It("can accept variable from env", func() {
			template := `---
files:
- name: test.yml
  content: {{ .variable }}`
			t := NewServiceTest()
			t.Environment().Set("CASTER_VAR_variable", "test")
			t.SetupString(template, &cast.Request{
				Template: "/template",
				Target:   "/output",
				Variables: []models.Variable{
					{
						Env: "CASTER_VAR_variable",
					},
				}})

			t.AssertExists("/output/test.yml")
			t.AssertContents("/output/test.yml", "test")
		})
	})
})
