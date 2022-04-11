package sr

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/flandiayingman/arkwaifu-2x/internal/app/dto"
	"github.com/flandiayingman/arkwaifu-2x/internal/app/util/dirutil"
	"github.com/flandiayingman/arkwaifu-2x/internal/app/util/executil"
	"github.com/flandiayingman/arkwaifu-2x/internal/app/util/pathutil"
	"github.com/flandiayingman/arkwaifu-2x/internal/app/vips"
	"go.uber.org/zap"
)

func init() {
	var err error
	switch runtime.GOOS {
	case "windows":
		err = executil.AddPath("./tools/realcugan-ncnn-vulkan-20220318-windows")
	case "linux":
		err = executil.AddPath("./tools/realcugan-ncnn-vulkan-20220318-ubuntu")
	default:
		log.Panic(fmt.Errorf("unsupported OS: %s", runtime.GOOS))
		return
	}
	if err != nil {
		zap.S().Errorw("Failed to add vips to PATH", "error", err)
		return
	}

	_, err = exec.LookPath(cuganExec)
	if err != nil {
		zap.S().Errorw("Failed to search for vips in PATH", "error", err)
		return
	}
}

type cugan struct {
	Name   string
	Params []string
}

const cuganExec = "realcugan-ncnn-vulkan"

func newCugan() (Model, error) {
	m := cugan{}
	m.Name = "real-cugan"
	m.Params = []string{
		// "-j", "1:1:1",
		"-s", "4",
	}
	log.Debugw("created real-cugan model",
		"model", m,
	)
	return &m, nil
}

func (m *cugan) Up(v dto.Variant, dir string) (dto.Variant, error) {
	srcV := v
	interV := srcV
	interV.Variant = m.Name
	interV.Filename = pathutil.ReplaceExt(v.Filename, ".inter.webp")
	destV := interV
	destV.Filename = pathutil.ReplaceExt(v.Filename, ".webp")

	srcPath := filepath.Join(dir, srcV.Path())
	interPath := filepath.Join(dir, interV.Path())
	destPath := filepath.Join(dir, destV.Path())

	err := dirutil.MkParentAll(interPath)
	if err != nil {
		return dto.Variant{}, err
	}
	err = dirutil.MkParentAll(destPath)
	if err != nil {
		return dto.Variant{}, err
	}

	params := append(m.Params, "-i", srcPath, "-o", interPath)
	cmd := exec.Command(cuganExec, params...)

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
