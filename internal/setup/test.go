package setup

import (
	"github.com/patrickhuber/caster/pkg/abstract/env"
	"github.com/patrickhuber/caster/pkg/abstract/fs"
	"github.com/patrickhuber/caster/pkg/cast"
	"github.com/patrickhuber/caster/pkg/console"
	"github.com/patrickhuber/caster/pkg/interpolate"
	"github.com/patrickhuber/go-di"
)

func NewTest() Setup {
	container := di.NewContainer()
	container.RegisterConstructor(env.NewMemory)
	container.RegisterConstructor(fs.NewMemory)
	container.RegisterConstructor(cast.NewService)
	container.RegisterConstructor(interpolate.NewService)
	container.RegisterConstructor(console.NewMemory)
	return &test{
		container: container,
	}
}

type test struct {
	container di.Container
}

func (t *test) Container() di.Container {
	return t.container
}

func (t *test) Close() error {
	return nil
}
