package helpers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gomig/utils"
)

// Download download file and save to destination
func Download(url string, dest string) error {
	// download file
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to download with status code: %d", resp.StatusCode)
	}

	// create dir
	if err := utils.CreateDirectory(filepath.Dir(dest)); err != nil {
		return err
	}

	// create the file
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	if _, err = io.Copy(out, resp.Body); err != nil {
		return err
	}

	return nil
}
