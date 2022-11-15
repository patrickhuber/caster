package interpolate_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/patrickhuber/caster/pkg/abstract/env"
	afs "github.com/patrickhuber/caster/pkg/abstract/fs"
	"github.com/patrickhuber/caster/pkg/interpolate"
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
`
			fs.Write("/template/.caster.yml", []byte(template), 0600)
			resp, err := svc.Interpolate(&interpolate.Request{
				Directory: "/template",
			})
			Expect(err).To(BeNil())
			Expect(resp).ToNot(BeNil())
		})
	})
})
