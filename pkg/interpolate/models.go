package interpolate

import "github.com/patrickhuber/caster/pkg/models"

// Request is the request object for casting a template
type Request struct {
	File      string            `yaml:"omitempty"`
	Directory string            `yaml:"omitempty"`
	Variables []models.Variable `yaml:"omitempty"`
}

type Response struct {
	SourceFile string        `yaml:"omitempty"`
	Caster     models.Caster `yaml:"omitempty"`
}
