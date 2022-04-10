package sr

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/flandiayingman/arkwaifu-2x/internal/app/dto"
	"github.com/flandiayingman/arkwaifu-2x/internal/app/util/executil"
	"github.com/flandiayingman/arkwaifu-2x/internal/app/util/pathutil"
	"github.com/flandiayingman/arkwaifu-2x/internal/app/vips"
	"go.uber.org/zap"
)

func init() {
	var err error
	switch runtime.GOOS {
	case "windows":
		err = executil.AddPath("./tools/realesrgan-ncnn-vulkan-20211212-windows")
	case "linux":
		err = executil.AddPath("./tools/realesrgan-ncnn-vulkan-20211212-ubuntu")
	default:
		log.Panic(fmt.Errorf("unsupported OS: %s", runtime.GOOS))
		return
	}
	if err != nil {
		zap.S().Errorw("Failed to add vips to PATH", "error", err)
		return
	}

	_, err = exec.LookPath(esrganExec)
	if err != nil {
		zap.S().Errorw("Failed to search for vips in PATH", "error", err)
		return
	}
}

type esrgan struct {
	Name   string
	Params []string
}

const esrganExec = "realesrgan-ncnn-vulkan"

func newEsrgan() (Model, error) {
	m := esrgan{}
	m.Name = "real-esrgan"
	m.Params = []string{
		// "-j", "1:1:1"
	}
	log.Debugw("created real-esrgan model",
		"model", m,
	)
	return &m, nil
}

func (m *esrgan) Up(v dto.Variant, dir string) (dto.Variant, error) {
	srcV := v
	interV := srcV
	interV.Variant = m.Name
	interV.Filename = pathutil.ReplaceExt(v.Filename, ".png")
	destV := interV
	destV.Filename = pathutil.ReplaceExt(v.Filename, ".webp")

	srcPath := filepath.Join(dir, srcV.Path())
	interPath := filepath.Join(dir, interV.Path())
	destPath := filepath.Join(dir, destV.Path())

	params := append(m.Params, "-i", srcPath, "-o", interPath)
	cmd := exec.Command(esrganExec, params...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return dto.Variant{}, err
	}
	err = executil.ScanOutput(m.Name, bytes.NewReader(output))
	if err != nil {
		return dto.Variant{}, err
	}
	defer func() { _ = os.Remove(interPath) }()

	err = vips.ConvertToWebp(interPath, destPath)
	if err != nil {
		return dto.Variant{}, err
	}

	return destV, nil
}
