package app

import (
	"context"
	"fmt"

	"github.com/flandiayingman/arkwaifu-2x/internal/app/dto"
	"github.com/flandiayingman/arkwaifu-2x/internal/app/sr"
	"github.com/pterm/pterm"
	"golang.org/x/sync/errgroup"
)

func UpscaleAssets(assets []dto.Asset) ([]dto.Variant, error) {
	progressbar, _ := progressbar.
		WithTitle("Upscaling the assets...").
		WithTotal(len(assets)).
		Start()
	vs := make([]dto.Variant, 0, len(assets))
	for _, a := range assets {
		progressbar.UpdateTitle(fmt.Sprintf("Upscaling the asset %s...", a.Name))
		v, err := upscaleParallel(a)
		if err != nil {
			pterm.Error.Printfln("Failed to upscale the asset %s.", a.Name)
			return nil, err
		}
		vs = append(vs, v...)
		pterm.Success.Printfln("Successfully upscaled the asset %s.", a.Name)
		progressbar.Increment()
	}
	return vs, nil
}

// 8 - 7m14s
func upscaleParallel(asset dto.Asset) ([]dto.Variant, error) {
	eg, _ := errgroup.WithContext(context.Background())
	vc := make(chan dto.Variant, len(asset.UpscalableVariants))
	for _, v := range asset.UpscalableVariants {
		v := v
		eg.Go(func() error {
			v, err := up(asset.Variants["img"], v)
			if err != nil {
				return err
			}
			vc <- v
			return nil
		})
	}
	err := eg.Wait()
	if err != nil {
		return nil, err
	}
	close(vc)

	var vs []dto.Variant
	for v := range vc {
		vs = append(vs, v)
	}
	return vs, nil
}

func up(srcv dto.Variant, v dto.Variant) (dto.Variant, error) {
	model := sr.Models[v.Variant]
	v, err := model.Up(srcv, dir)
	if err != nil {
		return dto.Variant{}, err
	}
	return v, nil
}
