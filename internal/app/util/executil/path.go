package executil

import (
	"fmt"
	"os"
	"path/filepath"
)

func AddPath(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	oldPath := os.Getenv("PATH")
	newPath := fmt.Sprintf("%s%s%s", absPath, string(os.PathListSeparator), oldPath)
	err = os.Setenv("PATH", newPath)
	if err != nil {
		return err
	}
	return nil
}
