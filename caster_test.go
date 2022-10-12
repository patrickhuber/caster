package main_test

import (
	"bytes"
	"text/template"

	"github.com/masterminds/sprig"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/caster/vfs"
)

var _ = Describe("Caster", func() {
	Describe("templates", func() {
		It("can render range of nested map", func() {
			t := template.New("test")
			t, err := t.Parse("{{range .regions}}{{.name}}{{end}}")
			Expect(err).To(BeNil())

			var writer bytes.Buffer
			data := map[string]interface{}{
				"regions": []map[string]interface{}{
					{"name": "eastus"},
					{"name": "westus"},
				},
			}
			err = t.Execute(&writer, data)
			Expect(err).To(BeNil())
			Expect(writer.String()).To(Equal("eastuswestus"))
		})
		It("can render sub type of nested map", func() {
			t := template.New("test")
			t, err := t.Parse("{{range .regions}}{{.name}}{{range .environments}}{{.name}}{{end}}{{end}}")
			Expect(err).To(BeNil())

			var writer bytes.Buffer
			data := map[string]interface{}{
				"regions": []map[string]interface{}{
					{
						"name": "eastus",
						"environments": []map[string]interface{}{
							{"name": "sdbx"},
							{"name": "nonprod"},
							{"name": "prod"},
						},
					},
					{
						"name": "westus",
						"environments": []map[string]interface{}{
							{"name": "sdbx"},
							{"name": "nonprod"},
							{"name": "prod"},
						},
					},
				},
			}
			err = t.Execute(&writer, data)
			Expect(err).To(BeNil())
			Expect(writer.String()).To(Equal("eastussdbxnonprodprodwestussdbxnonprodprod"))
		})
		It("can render newlines", func() {
			data := []string{"one", "two", "three"}
			t := template.New("test")
			t, err := t.Parse("{{range .}}{{.}}`n{{end}}")
			Expect(err).To(BeNil())

			var writer bytes.Buffer
			err = t.Execute(&writer, data)
			Expect(err).To(BeNil())
			Expect(writer.String()).To(Equal("one`ntwo`nthree`n"))
		})
		It("can render file from within template", func() {
			fs := vfs.NewMemory()
			err := fs.Write("test.yml", []byte("test: test"), 0600)
			Expect(err).To(BeNil())

			funcMap := map[string]interface{}{
				"file": func(path string) (string, error) {
					bytes, err := fs.Read(path)
					if err != nil {
						return "", err
					}
					return string(bytes), nil
				},
			}

			t, err := template.
				New("test").
				Funcs(funcMap).
				Parse("{{file \"test.yml\" }}")
			Expect(err).To(BeNil())

			var writer bytes.Buffer
			err = t.Execute(&writer, nil)
			Expect(err).To(BeNil())
			Expect(writer.String()).To(Equal("test: test"))
		})
		It("renders data after function", func() {
			fs := vfs.NewMemory()
			err := fs.Write("test.yml", []byte("test: {{ . }}"), 0600)
			Expect(err).To(BeNil())

			funcMap := map[string]interface{}{
				"templatefile": func(path string, data interface{}) (string, error) {
					content, err := fs.Read(path)
					if err != nil {
						return "", err
					}
					t, err := template.
						New("inner").
						Parse(string(content))
					if err != nil {
						return "", err
					}
					var writer bytes.Buffer
					err = t.Execute(&writer, data)
					return writer.String(), nil
				},
			}

			t, err := template.
				New("test").
				Funcs(funcMap).
				Parse("{{ templatefile \"test.yml\" .}}")
			Expect(err).To(BeNil())

			var writer bytes.Buffer
			err = t.Execute(&writer, "value")
			Expect(err).To(BeNil())
			Expect(writer.String()).To(Equal("test: value"))
		})

		It("sprig allows setting child objects", func() {
			fs := vfs.NewMemory()
			err := fs.Write("test.yml", []byte("- sub: {{ .sub.name }}\n  top: {{ .top}}"), 0600)
			Expect(err).To(BeNil())

			funcMap := sprig.TxtFuncMap()
			funcMap["templatefile"] = func(path string, data interface{}) (string, error) {
				content, err := fs.Read(path)
				if err != nil {
					return "", err
				}
				t, err := template.
					New("inner").
					Parse(string(content))
				if err != nil {
					return "", err
				}
				var writer bytes.Buffer
				err = t.Execute(&writer, data)
				return writer.String(), nil
			}
			t, err := template.
				New("test").
				Funcs(funcMap).
				Parse(`{{ range .test }}{{$params := dict "sub" . "top" $.top}}{{ templatefile "test.yml" $params }}
{{end}}`)
			Expect(err).To(BeNil())

			data := map[string]interface{}{
				"test": []map[string]interface{}{
					map[string]interface{}{
						"name": "one",
						"sub": []map[string]interface{}{
							map[string]interface{}{
								"name": "hi",
							},
						},
					},
					map[string]interface{}{
						"name": "two",
						"sub": []map[string]interface{}{
							map[string]interface{}{
								"name": "there",
							},
						},
					},
				},
				"top": "level",
			}
			var writer bytes.Buffer
			err = t.Execute(&writer, data)
			Expect(err).To(BeNil())
			Expect(writer.String()).To(Equal("- sub: one\n  top: level\n- sub: two\n  top: level\n"))
		})
	})
})
