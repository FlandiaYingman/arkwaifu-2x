package executil

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"go.uber.org/zap"
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

func ScanOutput(name string, output io.Reader) error {
	scanner := bufio.NewScanner(output)
	for scanner.Scan() {
		zap.S().Debugf("%s> %s", name, scanner.Text())
	}
	return scanner.Err()
}
