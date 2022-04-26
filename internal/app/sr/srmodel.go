package sr

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/flandiayingman/arkwaifu-2x/internal/app/dto"
	"github.com/flandiayingman/arkwaifu-2x/internal/app/util/fileutil"
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

var gpuLock sync.Mutex

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

func (m Model) Up(cv *dto.Variant, av *dto.Variant, dir string) (dto.Variant, error) {
	var err error
	dv := newDestinationVariant(*cv, m)
	cvPath := filepath.Join(dir, cv.Path())
	dvPath := filepath.Join(dir, dv.Path())

	done, _ := verify.Verify(dvPath)
	if done {
		return dv, nil
	}
	_ = fileutil.MkParents(dvPath)

	colorUpPath := filepath.Join(dir, cv.Path()+".color.png")
	defer func() { _ = os.Remove(colorUpPath) }()
	err = fileutil.MkParents(colorUpPath)
	if err != nil {
		return dto.Variant{}, err
	}
	err = m.up(cvPath, colorUpPath)
	if err != nil {
		return dto.Variant{}, err
	}

	if av != nil {
		avPath := filepath.Join(dir, av.Path())
		alphaUpPath := filepath.Join(dir, av.Path()+".alpha.png")
		defer func() { _ = os.Remove(alphaUpPath) }()
		err = fileutil.MkParents(alphaUpPath)
		if err != nil {
			return dto.Variant{}, err
		}
		err = m.up(avPath, alphaUpPath)
		if err != nil {
			return dto.Variant{}, err
		}

		err = vips.MergeAlpha(colorUpPath, alphaUpPath, dvPath)
		if err != nil {
			return dto.Variant{}, err
		}
	} else {
		err := vips.ConvertToWebp(colorUpPath, dvPath)
		if err != nil {
			return dto.Variant{}, err
		}
	}

	err = verify.Done(dvPath)
	if err != nil {
		return dto.Variant{}, err
	}

	return dv, nil
}

func (m Model) up(avPath string, alphaUpPath string) error {
	gpuLock.Lock()
	defer gpuLock.Unlock()
	err := m.upscale(avPath, alphaUpPath)
	if err != nil {
		return err
	}
	return nil
}

func newDestinationVariant(v dto.Variant, m model) dto.Variant {
	v.Variant = m.ModelName()
	v.Filename = pathutil.ReplaceExt(v.Filename, ".webp")
	return v
}
