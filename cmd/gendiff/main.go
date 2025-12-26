package main

import (
	"context"
	"os"

	"code/internal/app"
)

func main() {
	cmd := app.New()

	err := cmd.Run(context.Background(), os.Args)
	if err != nil {
		os.Exit(1)
	}
}
