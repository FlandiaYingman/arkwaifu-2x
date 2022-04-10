package sr

import (
	"github.com/flandiayingman/arkwaifu-2x/internal/app/dto"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

var log = zap.S()

type Model interface {
	Up(v dto.Variant, dir string) (dto.Variant, error)
}

func UpscalableVariants(variants []dto.Variant) []dto.Variant {
	var upscalable []dto.Variant

	for _, m := range ModelNames {
		// if "m" isn't in "variants", add a new variant to "upscalable"
		if !lo.ContainsBy(variants, func(v dto.Variant) bool { return v.Variant == m }) {
			upscalable = append(upscalable, dto.Variant{
				Kind:    variants[0].Kind,
				Name:    variants[0].Name,
				Variant: m,
			})
		}
	}

	return upscalable
}
