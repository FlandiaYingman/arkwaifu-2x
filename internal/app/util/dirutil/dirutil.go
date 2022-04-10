package dirutil

import (
	"os"
	"path/filepath"
)

func MkParentAll(path string) error {
	return os.MkdirAll(filepath.Dir(path), 0755)
}
