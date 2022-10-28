package caster

import (
	"os"

	"github.com/urfave/cli/v2"
)

// set with -ldflags
var version = ""

func main() {

	(&cli.App{}).Run(os.Args)
}
