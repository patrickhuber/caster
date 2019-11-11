package cast_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/caster/cast"
	"github.com/patrickhuber/caster/vfs"
)

var _ = Describe("Service", func() {
	Describe("Cast", func() {
		It("does not write caster file to output", func() {
			fs := vfs.NewMemory()
			err := fs.Write("/template/.caster", []byte(""), 600)
			Expect(err).To(BeNil())

			svc := cast.NewService(fs)
			err = svc.Cast("/template", "/output", nil)
			Expect(err).To(BeNil())

			// parent directory exists
			ok, err := fs.Exists("/output")
			Expect(err).To(BeNil())
			Expect(ok).To(BeTrue())

			// caster file doesn't exist
			ok, err = fs.Exists("/output/.caster")
			Expect(err).To(BeNil())
			Expect(ok).To(BeFalse())
		})
		It("writes plain files to target", func() {
			fs := vfs.NewMemory()
			err := fs.Write("/template/test.yml", []byte("test: test"), 600)
			Expect(err).To(BeNil())
			err = fs.Write("/template/sub/test.yml", []byte("test: sub"), 600)
			Expect(err).To(BeNil())

			svc := cast.NewService(fs)
			err = svc.Cast("/template", "/output", nil)
			Expect(err).To(BeNil())

			// parent directory exists
			ok, err := fs.Exists("/output")
			Expect(err).To(BeNil())
			Expect(ok).To(BeTrue())

			// test yml exists
			ok, err = fs.Exists("/output/test.yml")
			Expect(err).To(BeNil())
			Expect(ok).To(BeTrue())

			// contents are the same
			content, err := fs.Read("/output/test.yml")
			Expect(err).To(BeNil())
			Expect(string(content)).To(Equal("test: test"))

			// child directory exists
			ok, err = fs.Exists("/output/sub")
			Expect(err).To(BeNil())
			Expect(ok).To(BeTrue())

			// test yml exists
			ok, err = fs.Exists("/output/sub/test.yml")
			Expect(err).To(BeNil())
			Expect(ok).To(BeTrue())
		})
		It("evaluates file names", func() {
			fs := vfs.NewMemory()
			err := fs.Write("/template/{{&quot;hello&quot;}}{{&quot;world&quot;}}.yml", []byte("hello: world"), 0600)
			Expect(err).To(BeNil())

			svc := cast.NewService(fs)
			err = svc.Cast("/template", "/output", nil)
			Expect(err).To(BeNil())

			// parent directory exists
			ok, err := fs.Exists("/output")
			Expect(err).To(BeNil())
			Expect(ok).To(BeTrue())

			// test yml exists
			ok, err = fs.Exists("/output/helloworld.yml")
			Expect(err).To(BeNil())
			Expect(ok).To(BeTrue())
		})
	})
})
