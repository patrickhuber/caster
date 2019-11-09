package main_test

import (
	"bytes"
	"text/template"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
	})
})
