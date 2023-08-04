package setup

import (
	"github.com/patrickhuber/caster/internal/interpolate"
	"github.com/patrickhuber/go-xplat/filepath"
	"github.com/patrickhuber/go-xplat/fs"

	"github.com/patrickhuber/caster/internal/cast"
	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/go-xplat/env"

	"github.com/patrickhuber/go-xplat/console"
)

func New() Setup {
	container := di.NewContainer()
	container.RegisterConstructor(env.NewOS)
	container.RegisterConstructor(fs.NewOS)
	container.RegisterConstructor(func() *filepath.Processor {
		// options cause issues with constructor registration
		return filepath.NewProcessor()
	})
	container.RegisterConstructor(cast.NewService)
	container.RegisterConstructor(interpolate.NewService)
	container.RegisterConstructor(console.NewOS)
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
