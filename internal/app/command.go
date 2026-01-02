package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/urfave/cli/v3"

	"code/internal/diff"
	"code/internal/formatters"
	"code/internal/formatters/common"
	"code/internal/parser"
)

const (
	defaultFormat          = common.FormatStylish
	requiredArgumentsCount = 2

	formatFlag = "format"
)

func New() *cli.Command {
	return &cli.Command{
		Name:      "gendiff",
		Usage:     "Compares two configuration files and shows a difference.",
		ArgsUsage: "<filepath1> <filepath2>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    formatFlag,
				Aliases: []string{"f"},
				Value:   defaultFormat,
				Usage:   "output format",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() != requiredArgumentsCount {
				_ = cli.ShowRootCommandHelp(cmd)

				return cli.Exit(
					invalidArgsError(),
					exitCodeUsageError,
				)
			}

			firstFilepath := cmd.Args().Get(0)
			secondFilepath := cmd.Args().Get(1)
			outputFormat := cmd.String(formatFlag)

			output, err := Run(firstFilepath, secondFilepath, outputFormat)
			if err != nil {
				return cli.Exit(err.Error(), exitCodeFrom(err))
			}

			_, err = fmt.Fprintln(cmd.Writer, output)
			if err != nil {
				return err
			}

			return nil
		},
	}
}

func Run(filepath1, filepath2, format string) (string, error) {
	leftNode, err := parser.ParseFile(filepath1)
	if err != nil {
		return "", wrap(ErrRuntime, err)
	}

	rightNode, err := parser.ParseFile(filepath2)
	if err != nil {
		return "", wrap(ErrRuntime, err)
	}

	selectedFormatter, err := formatters.SelectFormatter(format)
	if err != nil {
		return "", wrap(ErrUsage, err)
	}

	changes := diff.Compare(leftNode, rightNode)

	return selectedFormatter.Format(changes), nil
}

func exitCodeFrom(err error) int {
	switch {
	case errors.Is(err, ErrUsage):
		return exitCodeUsageError
	default:
		return exitCodeRuntimeError
	}
}
