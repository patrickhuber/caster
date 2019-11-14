package cast

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"text/template"
	"text/template/parse"

	"github.com/patrickhuber/caster/vfs"
)

type generator struct {
	fs vfs.FileSystem
}

func (g *generator) Generate(root string, data map[string]interface{}) (string, error) {
	var buffer bytes.Buffer
	err := g.build(root, &buffer, 0)
	if err != nil {
		return "", err
	}
	content := buffer.String()
	tpl, err := template.New("").Parse(content)
	if err != nil {
		return "", err
	}

	var result bytes.Buffer
	err = tpl.Execute(&result, data)
	if err != nil {
		return "", err
	}
	return result.String(), nil
}

func (g *generator) build(root string, buffer io.Writer, level int) error {
	info, err := g.fs.Stat(root)
	if err != nil {
		return err
	}
	name := info.Name()

	if level > 0 {
		fmt.Fprintln(buffer)
	}
	fmt.Fprintf(buffer, "%s%s", strings.Repeat("\t", level), name)
	if info.IsDir() {
		fmt.Fprintf(buffer, "/")
	}

	files, err := g.fs.ReadDir(root)
	if err != nil {
		return err
	}
	for _, file := range files {
		err = g.build(g.fs.Join(root, file.Name()), buffer, level+1)
		if err != nil {
			return err
		}
	}

	if !strings.Contains(name, "{{") {
		return nil
	}

	_, err = parse.Parse("caster", name, "{{", "}}")
	if err == nil {
		return nil
	}
	// check if printing {{end}} resolves the error
	// if it doesn't, return the error
	_, err = parse.Parse("caster", name+"{{end}}", "{{", "}}")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(buffer)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(buffer, "%s%s", strings.Repeat("\t", level), "{{end}}")
	return err
}
