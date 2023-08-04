package setup

import (
	"github.com/patrickhuber/caster/internal/cast"
	"github.com/patrickhuber/caster/internal/initialize"
	"github.com/patrickhuber/caster/internal/interpolate"
	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/go-xplat/console"
	"github.com/patrickhuber/go-xplat/env"
	"github.com/patrickhuber/go-xplat/filepath"
	"github.com/patrickhuber/go-xplat/fs"
	"github.com/patrickhuber/go-xplat/os"
)

func NewTest() Setup {
	container := di.NewContainer()
	container.RegisterConstructor(env.NewMemory)
	container.RegisterConstructor(func() os.OS {
		return os.NewLinuxMock()
	})
	container.RegisterConstructor(func(processor *filepath.Processor) fs.FS {
		// options cause issues with constructor registration
		return fs.NewMemory(fs.WithProcessor(processor))
	})
	container.RegisterConstructor(func(o os.OS) *filepath.Processor {
		return filepath.NewProcessorWithOS(o)
	})
	container.RegisterConstructor(cast.NewService)
	container.RegisterConstructor(interpolate.NewService)
	container.RegisterConstructor(initialize.NewService)
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
