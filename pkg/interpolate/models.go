package interpolate

import "github.com/patrickhuber/caster/pkg/models"

type Variable struct {
	File  string
	Key   string
	Value string
	Env   string
}

// Request is the request object for casting a template
type Request struct {
	File      string
	Directory string
	Variables []Variable
}

type Response struct {
	SourceFile string
	Caster     models.Caster
}
