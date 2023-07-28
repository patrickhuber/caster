package commands

import (
	"github.com/patrickhuber/caster/internal/global"
	"github.com/patrickhuber/caster/internal/interpolate"
	"github.com/patrickhuber/caster/internal/models"
	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/go-xplat/console"
	"github.com/patrickhuber/go-xplat/env"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

const (
	InterpolateTemplateFlag = "interpolate"
	InterpolateNameFlag     = "name"
	InterpolateVarFlag      = "var"
	InterpolateVarFileFlag  = "var-file"
)

var Interpolate = &cli.Command{
	Name:        "interpolate",
	Description: "interpolates the specified template and outputs the result",
	Usage:       "interpolates the specified template and outputs the result",
	UsageText:   "caster apply [-t|--template <TEMPLATEDIR|TEMPLATEFILE>] [-n|--name <TEMPLATENAME>]",
	Action:      InterpolateAction,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    InterpolateTemplateFlag,
			Aliases: []string{"t"},
			Value:   ".",
		},
		&cli.StringFlag{
			Name:    InterpolateNameFlag,
			Aliases: []string{"n"},
		},
		&cli.StringSliceFlag{
			Name: InterpolateVarFlag,
		},
		&cli.StringSliceFlag{
			Name:      InterpolateVarFileFlag,
			TakesFile: true,
		},
	},
}

type InterpolateCommand struct {
	Options     InterpolateOptions
	Environment env.Environment     `inject:""`
	Service     interpolate.Service `inject:""`
	Console     console.Console     `inject:""`
}

type InterpolateOptions struct {
	Template  string
	Name      string
	Variables []models.Variable
}

func InterpolateAction(ctx *cli.Context) error {

	cmd := &InterpolateCommand{}
	resolver := ctx.App.Metadata[global.DependencyInjectionContainer].(di.Resolver)
	err := di.Inject(resolver, cmd)
	if err != nil {
		return err
	}

	variables, err := getFlagVariables(ctx)
	if err != nil {
		return err
	}

	envVariables, err := getEnvironmentVariables(cmd.Environment)
	if err != nil {
		return err
	}

	cmd.Options = InterpolateOptions{
		Template:  ctx.String(InterpolateTemplateFlag),
		Name:      ctx.String(InterpolateNameFlag),
		Variables: append(variables, envVariables...),
	}

	return cmd.Execute()
}

func (cmd *InterpolateCommand) Execute() error {
	variables := []models.Variable{}
	for _, v := range cmd.Options.Variables {
		variables = append(variables, models.Variable{
			File:  v.File,
			Key:   v.Key,
			Value: v.Value,
			Env:   v.Env,
		})
	}

	// create apply request
	request := &interpolate.Request{
		Template:  cmd.Options.Template,
		Variables: variables,
	}
	resp, err := cmd.Service.Interpolate(request)
	if err != nil {
		return err
	}
	encoder := yaml.NewEncoder(cmd.Console.Out())
	encoder.SetIndent(2)
	defer encoder.Close()
	return encoder.Encode(resp.Caster)
}
