package helpers

import (
	"fmt"
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

// Fatal error and clean output directory
func Fatal(err string, dist ...string) {
	fmt.Printf("Failed: %s\n", err)
	for _, dir := range dist {
		if os.Remove(dir) != nil {
			os.RemoveAll(dir)
		}
	}
	os.Exit(1)
}

// Handle fatal error if not nil
func Handle(err error, dist ...string) {
	if err != nil {
		Fatal(err.Error(), dist...)
	}
}

// HandleV fatal on error or return value
func HandleV[T any](v T, err error, dist ...string) T {
	if err != nil {
		Fatal(err.Error(), dist...)
	}
	return v
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
	return "/" + strings.TrimLeft(strings.NewReplacer(`\`, `/`).Replace(path), "/")
}

// IsPathOf check if is subdirectory of path
func IsPathOf(path, dir string) bool {
	return strings.HasPrefix(path, NormalizePath(dir)+"/")
}
