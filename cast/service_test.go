package cast_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/caster/cast"
	"github.com/patrickhuber/caster/vfs"
	"gopkg.in/yaml.v2"
)

type ServiceTest interface {
	Setup(template *cast.Caster, source, target string)
	SetupString(content, source, target string)
	SetupBytes(content []byte, source, target string)
	AssertExists(path string)
	AssertContents(path, content string)
	FileSystem() vfs.FileSystem
}

type serviceTest struct {
	fs vfs.FileSystem
}

func NewServiceTest() ServiceTest {
	return &serviceTest{
		fs: vfs.NewMemory(),
	}
}

func (t *serviceTest) Setup(template *cast.Caster, source, target string) {
	content, err := yaml.Marshal(template)
	Expect(err).To(BeNil())
	t.SetupBytes(content, source, target)
}

func (t *serviceTest) SetupString(content, source, target string) {
	t.SetupBytes([]byte(content), source, target)
}

func (t *serviceTest) SetupBytes(content []byte, source, target string) {
	templatePath := t.fs.Join(source, ".caster.yml")
	err := t.fs.Write(templatePath, content, 600)
	Expect(err).To(BeNil())

	svc := cast.NewService(t.fs)
	err = svc.Cast(source, target, nil)
	Expect(err).To(BeNil())

	// source directory exists
	t.AssertExists(target)
}

func (t *serviceTest) AssertExists(path string) {
	ok, err := t.fs.Exists(path)
	Expect(err).To(BeNil())
	Expect(ok).To(BeTrue(), "expected '%s' to exist", path)
}

func (t *serviceTest) AssertContents(path, content string) {
	data, err := t.fs.Read(path)
	Expect(err).To(BeNil())
	Expect(string(data)).To(Equal(content))
}

func (t *serviceTest) FileSystem() vfs.FileSystem {
	return t.fs
}

var _ = Describe("Service", func() {
	Describe("Cast", func() {
		It("writes plain files to target", func() {
			template := &cast.Caster{
				Files: []cast.File{
					cast.File{
						Name:    "test.yml",
						Content: "test: test",
					},
				},
				Folders: []cast.Folder{
					cast.Folder{
						Name: "sub",
						Files: []cast.File{
							cast.File{
								Name: "test.yml",
							},
						},
					},
				},
			}
			t := NewServiceTest()
			t.Setup(template, "/template", "/output")
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
			t.SetupString(template, "/template", "/output")
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
			t.SetupString(template, "/template", "/output")
			t.AssertExists("/output/hello")
			t.AssertExists("/output/hello/1.yml")
			t.AssertContents("/output/hello/1.yml", "one: 1")
		})
	})
})
