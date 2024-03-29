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
		Name:        "caster",
		Description: "a file and directory templating cli",
		Version:     version,
		Metadata: map[string]interface{}{
			global.DependencyInjectionContainer: runtime.Container(),
		},
		Commands: []*cli.Command{
			commands.Apply,
			commands.Interpolate,
			commands.Initialize,
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
