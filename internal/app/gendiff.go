package app

import (
	"context"

	"github.com/urfave/cli/v3"
)

// New application constructor
func New() *cli.Command {
	return &cli.Command{
		Name:      "gendiff",
		Usage:     "Compares two configuration files and shows a difference.",
		ArgsUsage: "",
		Action: func(_ context.Context, cmd *cli.Command) error {
			return nil
		},
	}
}
