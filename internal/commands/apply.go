package commands

import (
	"github.com/patrickhuber/caster/internal/global"
	"github.com/patrickhuber/caster/pkg/cast"
	"github.com/patrickhuber/go-di"
	"github.com/urfave/cli/v2"
)

const (
	ApplyFileFlag      = "apply"
	ApplyDirectoryFlag = "directory"
	ApplyNameFlag      = "name"
)

var Apply = &cli.Command{
	Name:        "apply",
	Description: "applies the specified template to the target directory",
	Usage:       "Applies the specified template to the target directory",
	UsageText:   "caster apply (-f <TEMPLATE_FILE> | -d <TEMPLATE_DIRECTORY | -n <NAME>) [<TARGET>]",
	Action:      ApplyAction,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    ApplyFileFlag,
			Aliases: []string{"f"},
		},
		&cli.StringFlag{
			Name:    ApplyDirectoryFlag,
			Aliases: []string{"d"},
		},
		&cli.StringFlag{
			Name:    ApplyNameFlag,
			Aliases: []string{"n"},
		},
	},
}

type ApplyCommand struct {
	Options ApplyOptions
	Service cast.Service
}

type ApplyOptions struct {
	Directory     string
	File          string
	Name          string
	Target        string
	VariableFiles []string
}

func (cmd *ApplyCommand) Execute() error {
	// create apply request
	request := &cast.CastRequest{
		Directory:     cmd.Options.Directory,
		File:          cmd.Options.File,
		VariableFiles: cmd.Options.VariableFiles,
		Target:        cmd.Options.Target,
	}
	err := cmd.Service.Cast(request)
	return err
}

func ApplyAction(ctx *cli.Context) error {
	resolver := ctx.App.Metadata[global.DependencyInjectionContainer].(di.Resolver)
	service, err := di.Resolve[cast.Service](resolver)
	if err != nil {
		return err
	}
	cmd := &ApplyCommand{
		Options: ApplyOptions{
			Directory: ctx.String(ApplyDirectoryFlag),
			File:      ctx.String(ApplyFileFlag),
			Name:      ctx.String(ApplyNameFlag),
			Target:    ctx.Args().First(),
		},
		Service: service,
	}

	return cmd.Execute()
}
