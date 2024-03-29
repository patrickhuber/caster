package commands

import (
	"fmt"
	"strings"

	"github.com/patrickhuber/caster/internal/cast"
	"github.com/patrickhuber/caster/internal/global"
	"github.com/patrickhuber/caster/internal/models"
	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/go-xplat/console"
	"github.com/patrickhuber/go-xplat/env"
	"github.com/urfave/cli/v2"
)

const (
	ApplyTemplateFlag = "template"
	ApplyNameFlag     = "name"
	ApplyOutFlag      = "out"
	ApplyVarFlag      = "var"
	ApplyVarFileFlag  = "var-file"
)

var Apply = &cli.Command{
	Name:        "apply",
	Description: "applies the specified template to the target directory",
	Usage:       "Applies the specified template to the target directory",
	UsageText:   "caster apply [-t|--template <TEMPLATEDIR|TEMPLATEFILE>] [-n|--name <TEMPLATENAME>] [OUTDIR]",
	Action:      ApplyAction,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    ApplyTemplateFlag,
			Aliases: []string{"t"},
			Value:   ".caster.yml",
		},
		&cli.StringFlag{
			Name:    ApplyNameFlag,
			Aliases: []string{"n"},
		},
		&cli.StringSliceFlag{
			Name: ApplyVarFlag,
		},
		&cli.StringSliceFlag{
			Name:      ApplyVarFileFlag,
			TakesFile: true,
		},
	},
}

type ApplyCommand struct {
	Options     ApplyOptions
	Environment env.Environment `inject:""`
	Service     cast.Service    `inject:""`
	Console     console.Console `inject:""`
}

type ApplyOptions struct {
	Template  string
	Name      string
	Target    string
	Variables []models.Variable
}

func (cmd *ApplyCommand) Execute() error {
	var variables []models.Variable

	// clone the variable slice
	variables = append(variables, cmd.Options.Variables...)

	// create apply request
	request := &cast.Request{
		Template:  cmd.Options.Template,
		Variables: variables,
		Target:    cmd.Options.Target,
	}
	err := cmd.Service.Cast(request)
	return err
}

func ApplyAction(ctx *cli.Context) error {
	cmd := &ApplyCommand{}
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

	cmd.Options = ApplyOptions{
		Template:  ctx.String(ApplyTemplateFlag),
		Name:      ctx.String(ApplyNameFlag),
		Target:    ctx.Args().First(),
		Variables: append(variables, envVariables...),
	}

	return cmd.Execute()
}

func getFlagVariables(ctx *cli.Context) ([]models.Variable, error) {
	variables := []models.Variable{}

	names := []string{}
	args := ctx.Args().Slice()
	for _, a := range args {
		if strings.Contains(a, "-"+ApplyVarFileFlag) {
			names = append(names, ApplyVarFileFlag)
		} else if strings.Contains(a, "-"+ApplyVarFlag) {
			names = append(names, ApplyVarFlag)
		}
	}
	varFlags := ctx.StringSlice(ApplyVarFlag)
	varFileFlags := ctx.StringSlice(ApplyVarFileFlag)
	varIndex := 0
	varFileIndex := 0
	for _, name := range names {
		switch name {
		case ApplyVarFlag:
			varFlag := varFlags[varIndex]
			split := strings.Split(varFlag, "=")
			if len(split) != 2 {
				return nil, fmt.Errorf("unable to parse var flag '%s'. Expected flag in format --var \"key=value\"", varFlag)
			}
			variables = append(variables, models.Variable{Key: split[0], Value: split[1]})
			varIndex++
		case ApplyVarFileFlag:
			variables = append(variables, models.Variable{File: varFileFlags[varFileIndex]})
			varFileIndex++
		}
	}
	return variables, nil
}

// getEnvironmentVariables returns the list of the environment variable keys that match the caster prefix
func getEnvironmentVariables(e env.Environment) ([]models.Variable, error) {
	variables := []models.Variable{}
	for k := range e.Export() {
		if !strings.HasPrefix(k, "CASTER_VAR_") {
			continue
		}
		variables = append(variables, models.Variable{Env: k})
	}
	return variables, nil
}
