package utils

import (
	"path/filepath"
	"strings"
)

func GetPathSegments(path string) []string {
	return strings.Split(path, "/")
}

func SplitPath(path string) (string, string) {
	return filepath.Split(path)
}
