package setup

import (
	"github.com/patrickhuber/caster/internal/cast"
	"github.com/patrickhuber/caster/internal/interpolate"
	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/go-xplat/console"
	"github.com/patrickhuber/go-xplat/env"
	"github.com/patrickhuber/go-xplat/filepath"
	"github.com/patrickhuber/go-xplat/fs"
	"github.com/patrickhuber/go-xplat/platform"
)

func NewTest() Setup {
	container := di.NewContainer()
	container.RegisterConstructor(env.NewMemory)
	container.RegisterConstructor(func(processor filepath.Processor) fs.FS {
		// options cause issues with constructor registration
		return fs.NewMemory(fs.WithProcessor(processor))
	})
	container.RegisterConstructor(func() filepath.Processor {
		return filepath.NewProcessorWithPlatform(platform.Linux)
	})
	container.RegisterConstructor(cast.NewService)
	container.RegisterConstructor(interpolate.NewService)
	container.RegisterConstructor(func() console.Console {
		return console.NewMemory()
	})
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
