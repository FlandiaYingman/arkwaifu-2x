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
		err = executil.AddPath("./tools/realcugan-ncnn-vulkan-20220728-windows")
	case "linux":
		err = executil.AddPath("./tools/realcugan-ncnn-vulkan-20220728-ubuntu")
	default:
		log.Panic(fmt.Errorf("unsupported OS: %s", runtime.GOOS))
		return
	}
	if err != nil {
		zap.S().Errorw("Failed to add cugan to PATH", "error", err)
		return
	}

	_, err = exec.LookPath(cuganExec)
	if err != nil {
		zap.S().Errorw("Failed to search for cugan in PATH", "error", err)
		return
	}
}

type cugan struct {
	Name   string
	Exec   string
	Params []string
}

const cuganExec = "realcugan-ncnn-vulkan"

func newCugan() (model, error) {
	m := cugan{
		Name: "real-cugan",
		Exec: cuganExec,
		Params: []string{
			// "-j", "1:1:1",
			"-s", "4",
		},
	}
	log.Debugw("created real-cugan model", "model", m)
	return &m, nil
}

func (m *cugan) ModelName() string {
	return m.Name
}

func (m *cugan) upscale(src, dst string) error {
	executable := m.Exec
	params := append(m.Params, "-i", src, "-o", dst)
	err := executil.Execute(m.ModelName(), executable, params...)
	if err != nil {
		return err
	}
	return nil
}
