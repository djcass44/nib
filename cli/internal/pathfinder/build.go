package pathfinder

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

var DefaultBuildPaths = []string{
	"dist",
	"build",
}

func FindBuildDir(dir string, buildPaths []string) (string, error) {
	var buildDir string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			return nil
		}
		if strings.Contains(path, "node_modules") {
			return nil
		}
		for _, p := range buildPaths {
			if d.Name() == p {
				buildDir = path
				return nil
			}
		}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("walking dir %s: %w", dir, err)
	}
	return buildDir, nil
}
