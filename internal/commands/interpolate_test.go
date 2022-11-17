package commands_test

import (
	"bytes"

	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/caster/internal/commands"
	"github.com/patrickhuber/caster/internal/global"
	"github.com/patrickhuber/caster/internal/setup"
	"github.com/patrickhuber/caster/pkg/abstract/env"
	"github.com/patrickhuber/caster/pkg/abstract/fs"
	"github.com/patrickhuber/caster/pkg/console"
	"github.com/patrickhuber/go-di"
	"github.com/urfave/cli/v2"
)

var _ = Describe("Interpolate", func() {
	var app *cli.App
	var container di.Container
	var con console.Console
	var f fs.FS
	var e env.Env
	BeforeEach(func() {
		var err error
		s := setup.NewTest()
		container = s.Container()

		con, err = di.Resolve[console.Console](container)
		Expect(err).To(BeNil())

		f, err = di.Resolve[fs.FS](container)
		Expect(err).To(BeNil())

		e, err = di.Resolve[env.Env](container)
		Expect(err).To(BeNil())

		app = &cli.App{
			Version: "1.0.0",
			Metadata: map[string]interface{}{
				global.DependencyInjectionContainer: container,
			},
			Commands: []*cli.Command{
				commands.Interpolate,
			},
			Reader:    con.In(),
			ErrWriter: con.Error(),
			Writer:    con.Out(),
		}

	})
	It("can run", func() {
		f.Write("/template/.caster.yml", []byte("files:\n- name: test.txt\n"), 0600)
		args := []string{"caster", "interpolate", "-d", "/template"}
		app.Metadata[global.OSArgs] = args

		err := app.Run(args)
		Expect(err).To(BeNil())

		buf, ok := con.Out().(*bytes.Buffer)
		Expect(ok).To(BeTrue())
		Expect(buf.String()).To(Equal("files:\n  - name: test.txt\n"))
	})
	When("env var", func() {
		It("can run", func() {
			template := `files:
  - name: test.txt
    content: {{ .key }}
`
			f.Write("/template/.caster.yml", []byte(template), 0600)
			e.Set("CASTER_VAR_key", "value")

			args := []string{"caster", "interpolate", "-d", "/template"}
			app.Metadata[global.OSArgs] = args

			err := app.Run(args)
			Expect(err).To(BeNil())

			buf, ok := con.Out().(*bytes.Buffer)
			Expect(ok).To(BeTrue())
			want := `files:
  - name: test.txt
    content: value
`
			have := buf.String()
			Expect(have).To(Equal(want), cmp.Diff(have, want))
		})
	})
	When("multiple data file", func() {
		It("can run", func() {
			f.Write("/template/.caster.yml", []byte("files:\n- name: test.txt\n  content: {{.first}}{{.second}}"), 0600)
			f.Write("/data/1.yml", []byte("first: first"), 0600)
			f.Write("/data/2.yml", []byte("second: second"), 0600)

			args := []string{"caster", "interpolate", "--var-file", "/data/1.yml", "--var-file", "/data/2.yml", "-d", "/template"}
			app.Metadata[global.OSArgs] = args
			err := app.Run(args)
			Expect(err).To(BeNil())

			buf, ok := con.Out().(*bytes.Buffer)
			Expect(ok).To(BeTrue())
			want := `files:
  - name: test.txt
    content: firstsecond
`
			have := buf.String()
			Expect(have).To(Equal(want), cmp.Diff(have, want))
		})
	})
	When("multiple argument", func() {
		It("can run", func() {
			f.Write("/template/.caster.yml", []byte("files:\n- name: test.txt\n  content: {{.first}}{{.second}}"), 0600)

			args := []string{"caster", "interpolate", "--var", "first=first", "--var", "second=second", "-d", "/template"}
			app.Metadata[global.OSArgs] = args
			err := app.Run(args)
			Expect(err).To(BeNil())

			buf, ok := con.Out().(*bytes.Buffer)
			Expect(ok).To(BeTrue())
			want := `files:
  - name: test.txt
    content: firstsecond
`
			have := buf.String()
			Expect(have).To(Equal(want), cmp.Diff(have, want))
		})
	})
	When("mixed arguments", func() {
		It("can override with var", func() {})
		It("can override with var-file", func() {})
	})
})
