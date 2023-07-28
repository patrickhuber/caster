package interpolate_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/patrickhuber/caster/internal/interpolate"
	"github.com/patrickhuber/caster/internal/models"
	"github.com/patrickhuber/go-xplat/env"
	"github.com/patrickhuber/go-xplat/filepath"
	afs "github.com/patrickhuber/go-xplat/fs"
	"github.com/patrickhuber/go-xplat/platform"
)

var _ = Describe("Service", func() {
	var (
		fs   afs.FS
		e    env.Environment
		svc  interpolate.Service
		path filepath.Processor
	)
	BeforeEach(func() {
		path = filepath.NewProcessorWithPlatform(platform.Linux)
		fs = afs.NewMemory(afs.WithProcessor(path))
		Expect(fs.Mkdir("/", 0600)).To(BeNil())
		Expect(fs.Mkdir("/template", 0600)).To(BeNil())
		e = env.NewMemory()
		svc = interpolate.NewService(fs, e, path)
	})
	Describe("Interpolate", func() {
		It("can interpolate", func() {
			template := `---
files:
- name: test.txt
  content: test
`
			err := fs.WriteFile("/template/.caster.yml", []byte(template), 0600)
			Expect(err).To(BeNil())
			resp, err := svc.Interpolate(&interpolate.Request{
				Template: "/template",
			})
			Expect(err).To(BeNil())
			Expect(resp).ToNot(BeNil())
			Expect(resp.Caster).ToNot(Equal(models.Caster{}))
			Expect(len(resp.Caster.Files)).To(Equal(1))
			file := resp.Caster.Files[0]
			Expect(file.Content).To(Equal("test"))
			Expect(file.Name).To(Equal("test.txt"))
		})
		It("can interpolate with data", func() {

			template := `---
files:
- name: test.txt
  content: {{ .key }}
`
			err := fs.WriteFile("/template/.caster.yml", []byte(template), 0600)
			Expect(err).To(BeNil())
			resp, err := svc.Interpolate(&interpolate.Request{
				Template: "/template",
				Variables: []models.Variable{
					{
						Key:   "key",
						Value: "value",
					},
				},
			})
			Expect(err).To(BeNil())
			Expect(resp).ToNot(BeNil())
			Expect(resp.Caster).ToNot(Equal(models.Caster{}))
			Expect(len(resp.Caster.Files)).To(Equal(1))
			file := resp.Caster.Files[0]
			Expect(file.Content).To(Equal("value"))
			Expect(file.Name).To(Equal("test.txt"))
		})
	})
})
