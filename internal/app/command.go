package app

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"

	"code/internal/parser"
)

func New() *cli.Command {
	return &cli.Command{
		Name:      "gendiff",
		Usage:     "Compares two configuration files and shows a difference.",
		ArgsUsage: "<filepath1> <filepath2>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "format",
				Aliases: []string{"f"},
				Value:   "stylish",
				Usage:   "output format",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() != 2 {
				_ = cli.ShowRootCommandHelp(cmd)

				return cli.Exit("\nexpected 2 arguments: <filepath1> <filepath2>", 2)
			}

			filepath1 := cmd.Args().Get(0)
			filepath2 := cmd.Args().Get(1)

			file1, err := parser.ParseFile(filepath1)
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			file2, err := parser.ParseFile(filepath2)
			if err != nil {
				return cli.Exit(err.Error(), 1)
			}

			fmt.Println(file1, file2)

			return nil
		},
	}
}
