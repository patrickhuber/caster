package interpolate_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/patrickhuber/caster/pkg/abstract/env"
	afs "github.com/patrickhuber/caster/pkg/abstract/fs"
	"github.com/patrickhuber/caster/pkg/interpolate"
	"github.com/patrickhuber/caster/pkg/models"
)

var _ = Describe("Service", func() {
	var (
		fs  afs.FS
		e   env.Env
		svc interpolate.Service
	)
	BeforeEach(func() {
		fs = afs.NewMemory()
		e = env.NewMemory()
		svc = interpolate.NewService(fs, e)
	})
	Describe("Interpolate", func() {
		It("can interpolate", func() {
			template := `---
files:
- name: test.txt
  content: test
`
			fs.Write("/template/.caster.yml", []byte(template), 0600)
			resp, err := svc.Interpolate(&interpolate.Request{
				Directory: "/template",
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
			fs.Write("/template/.caster.yml", []byte(template), 0600)
			resp, err := svc.Interpolate(&interpolate.Request{
				Directory: "/template",
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
