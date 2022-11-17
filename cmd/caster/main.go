package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/patrickhuber/caster/internal/commands"
	"github.com/patrickhuber/caster/internal/global"
	"github.com/patrickhuber/caster/internal/setup"
)

// set with -ldflags
var version = ""

func main() {
	runtime := setup.New()
	app := &cli.App{
		Version: version,
		Metadata: map[string]interface{}{
			global.DependencyInjectionContainer: runtime.Container(),
			global.OSArgs:                       os.Args,
		},
		Commands: []*cli.Command{
			commands.Apply,
			commands.Interpolate,
		},
	}
	err := app.Run(os.Args)
	handle(err)
}

func handle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
