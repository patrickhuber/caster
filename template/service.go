package template

// Service handles template lifecycle
type Service interface {
}

type service struct {
}

// NewService creates a new instance of the template service
func NewService() Service {
	return &service{}
}
