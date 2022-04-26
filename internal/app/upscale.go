package app

import (
	"github.com/flandiayingman/arkwaifu-2x/internal/app/dto"
	"github.com/flandiayingman/arkwaifu-2x/internal/app/sr"
	"go.uber.org/zap"
)

type UpscaleTask struct {
	colorVariant *dto.Variant
	alphaVariant *dto.Variant
	dstVariant   []*dto.Variant
	dstDir       *string
}

func (t *UpscaleTask) Run() {
	for _, v := range t.dstVariant {
		err := up(t.colorVariant, t.alphaVariant, v, *t.dstDir)
		if err != nil {
			zap.S().Warnw("Failed to upscale variant",
				"variant", v,
				"error", err,
			)
		}
		zap.S().Infow("Upscaled variant",
			"variant", v,
		)
	}
}
func up(cv, av, dv *dto.Variant, dstDir string) error {
	model := sr.Models[dv.Variant]
	v, err := model.Up(cv, av, dstDir)
	*dv = v
	return err
}

func (t *UpscaleTask) ToSubmitTask() *SubmitTask {
	return &SubmitTask{
		Variants: t.dstVariant,
		dstDir:   t.dstDir,
	}
}
