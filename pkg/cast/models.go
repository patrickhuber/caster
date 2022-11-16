package cast

import "github.com/patrickhuber/caster/pkg/models"

// Request is the request object for casting a template
type Request struct {
	File      string
	Directory string
	Target    string
	Variables []models.Variable
}
