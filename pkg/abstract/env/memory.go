package env

type memory struct {
	data map[string]string
}

func NewMemory() Env {
	return NewMemoryWith(map[string]string{})
}

func NewMemoryWith(data map[string]string) Env {
	return &memory{
		data: data,
	}
}

func (e *memory) Get(key string) string {
	return e.data[key]
}

func (e *memory) Lookup(key string) (string, bool) {
	ok, value := e.data[key]
	return ok, value
}

func (e *memory) Set(key, value string) error {
	e.data[key] = value
	return nil
}

func (e *memory) UnSet(key string) error {
	delete(e.data, key)
	return nil
}

func (e *memory) List() []string {
	list := []string{}
	for k := range e.data {
		list = append(list, k)
	}
	return list
}

func (e *memory) Clear() {
	for k := range e.data {
		delete(e.data, k)
	}
}
