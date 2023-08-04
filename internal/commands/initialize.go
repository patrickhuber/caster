package commands

import (
	"github.com/patrickhuber/caster/internal/global"
	"github.com/patrickhuber/caster/internal/initialize"
	"github.com/patrickhuber/go-di"
	"github.com/urfave/cli/v2"
)

const (
	InitializeTemplateFlag = "template"
)

var Initialize = &cli.Command{
	Name:        "initialize",
	Aliases:     []string{"init"},
	Description: "initializes the speified directory or file with the default template",
	Usage:       "initializes the speified directory or file with the default template",
	UsageText:   "caster init [DIRECTORY|FIILE]",
	Action:      InitializeAction,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    InitializeTemplateFlag,
			Aliases: []string{"t"},
			Value:   ".",
		},
	},
}

type InitializeCommand struct {
	Options InitializeOptions
	Service initialize.Service `inject:""`
}

type InitializeOptions struct {
	Template string
}

func InitializeAction(ctx *cli.Context) error {
	cmd := &InitializeCommand{}
	resolver := ctx.App.Metadata[global.DependencyInjectionContainer].(di.Resolver)
	err := di.Inject(resolver, cmd)
	if err != nil {
		return err
	}
	return cmd.Execute()
}

func (cmd *InitializeCommand) Execute() error {
	request := &initialize.Request{
		Template: cmd.Options.Template,
	}
	_, err := cmd.Service.Initialize(request)
	return err
}
