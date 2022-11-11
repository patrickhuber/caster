package cast_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/caster/pkg/cast"
	"gopkg.in/yaml.v3"

	afs "github.com/patrickhuber/caster/pkg/abstract/fs"
)

type ServiceTest interface {
	Setup(template *cast.Caster, source, target string)
	SetupWithFile(template *cast.Caster, sourceFile, target string)
	SetupString(content, source, target string)
	SetupBytes(content []byte, source, target string)
	AssertExists(path string)
	AssertContents(path, content string)
	FileSystem() afs.FS
}

type serviceTest struct {
	fs afs.FS
}

func NewServiceTest() ServiceTest {
	return &serviceTest{
		fs: afs.NewMemory(),
	}
}

func (t *serviceTest) Setup(template *cast.Caster, source, target string) {
	content, err := yaml.Marshal(template)
	Expect(err).To(BeNil())
	t.SetupBytes(content, source, target)
}

func (t *serviceTest) SetupWithFile(template *cast.Caster, sourceFile, target string) {
	content, err := yaml.Marshal(template)
	Expect(err).To(BeNil())
	t.SetupBytesWithFile(content, sourceFile, target)
}

func (t *serviceTest) SetupString(content, source, target string) {
	t.SetupBytes([]byte(content), source, target)
}

func (t *serviceTest) SetupBytes(content []byte, source, target string) {
	templatePath := t.fs.Join(source, ".caster.yml")
	t.SetupBytesWithFile(content, templatePath, target)
}

func (t *serviceTest) SetupBytesWithFile(content []byte, sourceFile, target string) {
	err := t.fs.Write(sourceFile, content, 0600)
	Expect(err).To(BeNil())

	source := t.fs.Dir(sourceFile)
	Expect(t.fs.IsDir(source)).To(BeTrue())

	svc := cast.NewService(t.fs)

	req := &cast.CastRequest{
		File:   sourceFile,
		Target: target,
	}
	err = svc.Cast(req)
	Expect(err).To(BeNil())

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

func (t *serviceTest) FileSystem() afs.FS {
	return t.fs
}

var _ = Describe("Service", func() {
	Describe("Cast", func() {
		When("caster file specified", func() {
			It("applies from specified file", func() {
				template := &cast.Caster{
					Files: []cast.File{
						{
							Name:    "test.yml",
							Content: "test: test",
						},
					},
					Folders: []cast.Folder{
						{
							Name: "sub",
							Files: []cast.File{
								{
									Name: "test.yml",
								},
							},
						},
					},
				}
				t := NewServiceTest()
				t.SetupWithFile(template, "/template/custom.yml", "/output")

				t.AssertExists("/output/test.yml")
				t.AssertContents("/output/test.yml", "test: test")
				t.AssertExists("/output/sub")
				t.AssertExists("/output/sub/test.yml")
			})
		})
		It("writes plain files to target", func() {
			template := &cast.Caster{
				Files: []cast.File{
					{
						Name:    "test.yml",
						Content: "test: test",
					},
				},
				Folders: []cast.Folder{
					{
						Name: "sub",
						Files: []cast.File{
							{
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
		It("writes ref", func() {
			template := `---
files:
- name: test
  ref: test.txt`
			t := NewServiceTest()

			err := t.FileSystem().Write("/template/test.txt", []byte("test"), 0644)
			Expect(err).To(BeNil())

			t.SetupString(template, "/template", "/output")
			t.AssertExists("/output/test")
			t.AssertContents("/output/test", "test")
		})
	})
})
