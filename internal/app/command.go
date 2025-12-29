package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/urfave/cli/v3"

	"code"
	"code/internal/formatter"
	"code/internal/formatter/stylish"
	"code/internal/parser"
)

const (
	formatFlag = "format"
	emptyFlag  = ""
)

const (
	defaultFormat          = formatter.FormatStylish
	requiredArgumentsCount = 2
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

			output, err := run(firstFilepath, secondFilepath, outputFormat)
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

func run(filepath1, filepath2, format string) (string, error) {
	leftNode, err := parser.ParseFile(filepath1)
	if err != nil {
		return "", wrap(ErrRuntime, err)
	}

	rightNode, err := parser.ParseFile(filepath2)
	if err != nil {
		return "", wrap(ErrRuntime, err)
	}

	selectedFormatter, err := selectFormatter(format)
	if err != nil {
		return "", wrap(ErrUsage, err)
	}

	return code.GenDiff(leftNode, rightNode, selectedFormatter), nil
}

func selectFormatter(format string) (formatter.Formatter, error) {
	switch format {
	case emptyFlag, formatter.FormatStylish:
		return stylish.New(), nil
	default:
		return nil, fmt.Errorf("%w: %s", formatter.ErrUnknownFormat, format)
	}
}

func exitCodeFrom(err error) int {
	switch {
	case errors.Is(err, ErrUsage):
		return exitCodeUsageError
	default:
		return exitCodeRuntimeError
	}
}
