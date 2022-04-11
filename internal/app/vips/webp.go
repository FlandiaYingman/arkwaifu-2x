package vips

import (
	"os"
	"os/exec"

	"github.com/flandiayingman/arkwaifu-2x/internal/app/util/executil"
	"go.uber.org/zap"
)

func init() {
	err := executil.AddPath("tools/vips-dev-8.12/bin")
	if err != nil {
		zap.S().Errorw("Failed to add vips to PATH", "error", err)
		return
	}
	_, err = exec.LookPath("vips")
	if err != nil {
		zap.S().Errorw("Failed to search for vips in PATH", "error", err)
		return
	}
}

const (
	vipsExec = "vips"
)

func ConvertToWebp(srcPath string, dstPath string) error {
	// See:
	// https://www.libvips.org/API/current/using-cli.html
	// https://www.libvips.org/API/current/VipsForeignSave.html
	format := "[Q=100,preset=VIPS_FOREIGN_WEBP_PRESET_PICTURE,strip]"
	err := executil.Execute("vips", vipsExec, "copy", srcPath, dstPath+format)
	if err != nil {
		return err
	}
	err = os.Remove(srcPath)
	if err != nil {
		return err
	}
	return nil
}
