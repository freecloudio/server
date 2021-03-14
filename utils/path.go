package utils

import (
	"path/filepath"
	"strings"
)

func GetPathSegments(path string) []string {
	allSegs := strings.Split(path, "/")
	filteredSegs := []string{}
	for _, seg := range allSegs {
		trimmedSeg := strings.TrimSpace(seg)
		if trimmedSeg != "" {
			filteredSegs = append(filteredSegs, trimmedSeg)
		}
	}
	return filteredSegs
}

func SplitPath(path string) (string, string) {
	return filepath.Split(path)
}
