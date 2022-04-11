package executil

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"

	"go.uber.org/zap"
)

func Execute(name string, executable string, params ...string) error {
	cmd := exec.Command(executable, params...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to run %s(%v): %w\noutput: \n%v", name, cmd, err, string(output))
	}
	err = scanOutput(name, bytes.NewReader(output))
	if err != nil {
		return fmt.Errorf("failed to scan output of %s(%v): %w\noutput: \n%v", name, cmd, err, string(output))
	}
	return nil
}

func scanOutput(name string, output io.Reader) error {
	scanner := bufio.NewScanner(output)
	zap.S().Debugf("### Begin %s ###", name)
	for scanner.Scan() {
		zap.S().Debugf("%s> %s", name, scanner.Text())
	}
	zap.S().Debugf("### End %s ###", name)
	return scanner.Err()
}
