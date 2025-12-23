package main

import (
	"context"
	"os"

	"code/internal/app"
)

func main() {
	cmd := app.New()

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		os.Exit(1)
	}
}
