package sr

import (
	"os"
	"path/filepath"

	"github.com/flandiayingman/arkwaifu-2x/internal/app/dto"
	"github.com/flandiayingman/arkwaifu-2x/internal/app/util/dirutil"
	"github.com/flandiayingman/arkwaifu-2x/internal/app/util/pathutil"
	"github.com/flandiayingman/arkwaifu-2x/internal/app/verify"
	"github.com/flandiayingman/arkwaifu-2x/internal/app/vips"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

var log = zap.S()

type Model struct {
	model
}

type model interface {
	ModelName() string
	upscale(src, dst string) error
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

func (m Model) Up(sv dto.Variant, dir string) (dto.Variant, error) {
	iv := newIntermediateVariant(sv, m)
	dv := newDestinationVariant(sv, m)
	svPath := filepath.Join(dir, sv.Path())
	ivPath := filepath.Join(dir, iv.Path())
	dvPath := filepath.Join(dir, dv.Path())
	defer func() { _ = os.Remove(ivPath) }()
	err := dirutil.MkParentAll(ivPath)
	if err != nil {
		return dto.Variant{}, err
	}
	err = dirutil.MkParentAll(dvPath)
	if err != nil {
		return dto.Variant{}, err
	}

	done, _ := verify.Verify(dvPath)
	if done {
		return dv, nil
	}

	err = m.upscale(svPath, ivPath)
	if err != nil {
		return dto.Variant{}, err
	}
	err = vips.ConvertToWebp(ivPath, dvPath)
	if err != nil {
		return dto.Variant{}, err
	}

	err = verify.Done(dvPath)
	if err != nil {
		return dto.Variant{}, err
	}

	return dv, nil
}

func newIntermediateVariant(v dto.Variant, m model) dto.Variant {
	v.Variant = m.ModelName()
	v.Filename = pathutil.ReplaceExt(v.Filename, ".inter.webp")
	return v
}
func newDestinationVariant(v dto.Variant, m model) dto.Variant {
	v.Variant = m.ModelName()
	v.Filename = pathutil.ReplaceExt(v.Filename, ".webp")
	return v
}
