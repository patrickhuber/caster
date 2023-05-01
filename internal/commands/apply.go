package commands

import (
	"fmt"
	"strings"

	"github.com/patrickhuber/caster/internal/global"
	"github.com/patrickhuber/caster/pkg/cast"
	"github.com/patrickhuber/caster/pkg/models"
	"github.com/patrickhuber/go-di"
	"github.com/patrickhuber/go-xplat/env"
	"github.com/urfave/cli/v2"
)

const (
	ApplyFileFlag      = "apply"
	ApplyDirectoryFlag = "directory"
	ApplyNameFlag      = "name"
	ApplyVarFlag       = "var"
	ApplyVarFileFlag   = "var-file"
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
}

type ApplyOptions struct {
	Directory string
	File      string
	Name      string
	Target    string
	Variables []models.Variable
}

func (cmd *ApplyCommand) Execute() error {
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
	request := &cast.Request{
		Directory: cmd.Options.Directory,
		File:      cmd.Options.File,
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
		Directory: ctx.String(ApplyDirectoryFlag),
		File:      ctx.String(ApplyFileFlag),
		Name:      ctx.String(ApplyNameFlag),
		Target:    ctx.Args().First(),
		Variables: append(variables, envVariables...),
	}

	return cmd.Execute()
}

func getFlagVariables(ctx *cli.Context) ([]models.Variable, error) {
	variables := []models.Variable{}

	names := []string{}
	args, _ := ctx.App.Metadata[global.OSArgs].([]string)
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
