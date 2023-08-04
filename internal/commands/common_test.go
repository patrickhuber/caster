package commands_test

import (
	"testing"

	"github.com/patrickhuber/caster/internal/commands"
	"github.com/patrickhuber/caster/internal/global"
	"github.com/patrickhuber/caster/internal/setup"
	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/go-xplat/console"
	"github.com/patrickhuber/go-xplat/env"
	"github.com/patrickhuber/go-xplat/filepath"
	"github.com/patrickhuber/go-xplat/fs"
	"github.com/patrickhuber/go-xplat/os"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

type TestContext struct {
	app       *cli.App
	container di.Container
	console   console.Console
	fs        fs.FS
	env       env.Environment
	os        os.OS
	path      *filepath.Processor
}

func SetupTestContext(t *testing.T) *TestContext {
	var err error
	s := setup.NewTest()
	container := s.Container()

	o, err := di.Resolve[os.OS](container)
	require.NoError(t, err)

	wd, err := o.WorkingDirectory()
	require.NoError(t, err)

	con, err := di.Resolve[console.Console](container)
	require.NoError(t, err)

	f, err := di.Resolve[fs.FS](container)
	require.NoError(t, err)

	paths := []string{"/", "/template", "/data", wd, o.Home()}
	for _, p := range paths {
		err = f.MkdirAll(p, 0666)
		require.NoError(t, err)
	}

	e, err := di.Resolve[env.Environment](container)
	require.NoError(t, err)

	p, err := di.Resolve[*filepath.Processor](container)
	require.NoError(t, err)

	app := &cli.App{
		Version: "1.0.0",
		Metadata: map[string]interface{}{
			global.DependencyInjectionContainer: container,
		},
		Commands: []*cli.Command{
			commands.Interpolate,
			commands.Apply,
			commands.Initialize,
		},
		Reader:    con.In(),
		ErrWriter: con.Error(),
		Writer:    con.Out(),
	}

	return &TestContext{
		app:       app,
		container: container,
		console:   con,
		fs:        f,
		env:       e,
		os:        o,
		path:      p,
	}
}
