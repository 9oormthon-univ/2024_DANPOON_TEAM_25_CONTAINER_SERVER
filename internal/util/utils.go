package util

import (
	"fmt"
	"os"
	"path/filepath"
)

func SetWorkingDirectory() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("project root with go.mod not found")
		}
		dir = parent
	}
}

func GetPath(relativePath string) string {
	projectRoot, _ := SetWorkingDirectory()

	filePath := filepath.Join(projectRoot, relativePath)
	return filePath
}
