package pathfinder

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"slices"
	"strings"
)

var DefaultBuildPaths = []string{
	"dist",
	"build",
}

var IgnoredFilePaths = []string{
	"node_modules",
}

// FindBuildDir attempts to locate the directory containing packaged static files.
// Generally this is in the DefaultBuildPaths list, however custom build frameworks
// may use different directories (e.g. public)
func FindBuildDir(dir string, buildPaths []string) (string, error) {
	var buildDir string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			return nil
		}
		// if we ever see the node_modules directory, back out.
		// We want to make sure we never accidentally package it in
		// our container
		if slices.ContainsFunc(IgnoredFilePaths, func(s string) bool { return strings.Contains(path, s) }) {
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
