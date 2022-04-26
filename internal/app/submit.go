package app

import (
	"github.com/flandiayingman/arkwaifu-2x/internal/app/api"
	"github.com/flandiayingman/arkwaifu-2x/internal/app/dto"
	"go.uber.org/zap"
)

type SubmitTask struct {
	Variants []*dto.Variant
	dstDir   *string
}

func (t *SubmitTask) Run() {
	err := api.PostVariants(t.Variants, *t.dstDir)
	if err != nil {
		zap.S().Warnw("Failed to submit the upscaled variants.",
			"variants", t.Variants,
			"error", err,
		)
	} else {
		zap.S().Infow("Submitted the upscaled variants.",
			"variants", t.Variants,
		)
	}
}
