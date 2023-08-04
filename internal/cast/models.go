package cast

import "github.com/patrickhuber/caster/internal/models"

// Request is the request object for casting a template
type Request struct {
	Template  string
	Target    string
	Variables []models.Variable
}
