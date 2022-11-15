package cast

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
	Target    string
	Variables []Variable
}
