package app

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"

	"code"
	"code/internal/formatter"
	"code/internal/formatter/stylish"
	"code/internal/parser"
)

const (
	exitCodeUsageError   = 2
	exitCodeRuntimeError = 1
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
					fmt.Sprintf(
						"\nexpected %d arguments: <filepath1> <filepath2>",
						requiredArgumentsCount,
					),
					exitCodeUsageError,
				)
			}

			firstFilepath := cmd.Args().Get(0)
			secondFilepath := cmd.Args().Get(1)

			outputFormat := cmd.String(formatFlag)

			leftNode, err := parser.ParseFile(firstFilepath)
			if err != nil {
				return cli.Exit(err.Error(), exitCodeRuntimeError)
			}

			rightNode, err := parser.ParseFile(secondFilepath)
			if err != nil {
				return cli.Exit(err.Error(), exitCodeRuntimeError)
			}

			selectedFormatter, err := selectFormatter(outputFormat)
			if err != nil {
				return cli.Exit(err.Error(), exitCodeUsageError)
			}

			output := code.GenDiff(leftNode, rightNode, selectedFormatter)

			fmt.Println(output)

			return nil
		},
	}
}

func selectFormatter(format string) (formatter.Formatter, error) {
	switch format {
	case emptyFlag, formatter.FormatStylish:
		return stylish.New(), nil
	default:
		return nil, fmt.Errorf("%w: %s", formatter.ErrUnknownFormat, format)
	}
}
