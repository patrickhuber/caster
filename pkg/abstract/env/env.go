package env

import (
	"os"
)

type Env interface {
	Get(key string) string
	Lookup(name string) (string, bool)
	Set(key, value string) error
	UnSet(key string) error
	List() []string
	Clear()
}

type env struct {
}

func NewOsEnv() Env {
	return &env{}
}

func (e *env) Get(key string) string {
	return os.Getenv(key)
}

func (e *env) Lookup(key string) (string, bool) {
	return os.LookupEnv(key)
}

func (e *env) Set(key, value string) error {
	return os.Setenv(key, value)
}

func (e *env) UnSet(key string) error {
	return os.Unsetenv(key)
}

func (e *env) List() []string {
	return os.Environ()
}

func (e *env) Clear() {
	os.Clearenv()
}
