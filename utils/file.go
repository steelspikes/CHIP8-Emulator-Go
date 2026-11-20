package utils

import (
	"os"
	"path/filepath"
)

func ReadFile(name string) ([]byte, error) {
	return os.ReadFile(filepath.Join("bin", name))
}
