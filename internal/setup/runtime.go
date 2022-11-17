package setup

import (
	"github.com/patrickhuber/caster/pkg/abstract/fs"
	"github.com/patrickhuber/caster/pkg/console"
	"github.com/patrickhuber/caster/pkg/interpolate"

	"github.com/patrickhuber/caster/pkg/abstract/env"
	"github.com/patrickhuber/caster/pkg/cast"
	"github.com/patrickhuber/go-di"
)

func New() Setup {
	container := di.NewContainer()
	container.RegisterConstructor(env.NewOsEnv)
	container.RegisterConstructor(fs.NewOs)
	container.RegisterConstructor(cast.NewService)
	container.RegisterConstructor(interpolate.NewService)
	container.RegisterConstructor(console.New)
	return &runtime{
		container: container,
	}
}

type runtime struct {
	container di.Container
}

func (r *runtime) Close() error {
	return nil
}

func (r *runtime) Container() di.Container {

	return r.container
}
