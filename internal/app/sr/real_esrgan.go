package sr

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/flandiayingman/arkwaifu-2x/internal/app/util/executil"
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
		zap.S().Errorw("Failed to add esrgan to PATH", "error", err)
		return
	}

	_, err = exec.LookPath(esrganExec)
	if err != nil {
		zap.S().Errorw("Failed to search for esrgan in PATH", "error", err)
		return
	}
}

type esrgan struct {
	Name   string
	Exec   string
	Params []string
}

const esrganExec = "realesrgan-ncnn-vulkan"

func newEsrgan() (model, error) {
	m := esrgan{
		Name:   "real-esrgan",
		Exec:   esrganExec,
		Params: []string{
			// "-j", "1:1:1",
		},
	}
	log.Debugw("created real-esrgan model", "model", m)
	return &m, nil
}

func (m *esrgan) ModelName() string {
	return m.Name
}

func (m *esrgan) upscale(src, dst string) error {
	executable := m.Exec
	params := append(m.Params, "-i", src, "-o", dst)
	err := executil.Execute(m.ModelName(), executable, params...)
	if err != nil {
		return err
	}
	return nil
}
