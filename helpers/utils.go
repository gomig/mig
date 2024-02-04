package helpers

import (
	"os"
	"strings"
)

// ReadFile read file data
func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// WriteFile write data to file
func WriteFile(path string, data []byte) error {
	return os.WriteFile(path, data, os.ModePerm)
}

// NormalizeRepo generate normalized github repo path
func NormalizeRepo(r string) string {
	r = strings.Replace(r, "https://github.com", "", 1)
	r = strings.Replace(r, "http://github.com", "", 1)
	r = strings.Replace(r, "github.com", "", 1)
	r = strings.TrimLeft(r, "/")
	return "https://github.com/" + r
}

// NormalizePath set path separator to /
func NormalizePath(path string) string {
	return "/" + strings.TrimRight(strings.TrimLeft(strings.NewReplacer(`\`, `/`).Replace(path), "/"), "/")
}

// IsPathOf check if is subdirectory of path
func IsPathOf(path, dir string) bool {
	path = NormalizePath(path)
	dir = NormalizePath(dir)
	return path == dir || strings.HasPrefix(path, dir+"/")
}
