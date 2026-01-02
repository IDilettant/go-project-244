package code

import "code/internal/app"

// GenDiff reads two config files, computes a diff, and renders it in the requested output format
func GenDiff(filepath1, filepath2, format string) (string, error) {
	return app.Run(filepath1, filepath2, format)
}
