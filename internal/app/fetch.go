package app

import (
	"math/rand"
	"time"

	"github.com/flandiayingman/arkwaifu-2x/internal/app/api"
	"github.com/flandiayingman/arkwaifu-2x/internal/app/dto"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type FetchTask struct {
	total     int
	remaining int
	asset     dto.Asset
	dstDir    *string
}

func CreateTask(dstDir *string) *FetchTask {
	assets, upscalableAssets, err := api.GetAssetsUpscalable()
	if err != nil {
		zap.S().Warnw("Failed to fetch assets.",
			"error", err,
		)
		return nil
	}
	if len(upscalableAssets) == 0 {
		zap.S().Infof("No upscalable assets found, remaining/total: %d/%d", len(upscalableAssets), len(assets))
		return nil
	}

	preference := func(a dto.Asset) bool { return a.Kind == "images" || a.Kind == "backgrounds" }
	upAssetsPrefers := lo.Filter(upscalableAssets, func(a dto.Asset, _ int) bool { return preference(a) })
	upAssetsOthers := lo.Filter(upscalableAssets, func(a dto.Asset, _ int) bool { return !preference(a) })

	// Select asset randomly to avoid 'race' to some degree.
	// The 'race' refers to, like, multiple clients are upscaling the same asset,
	// though at last they will get the same result, it's just a waste of time.
	var a dto.Asset
	if len(upAssetsPrefers) > 0 {
		a = upAssetsPrefers[rand.Intn(len(upAssetsPrefers))]
	} else {
		a = upAssetsOthers[rand.Intn(len(upAssetsOthers))]
	}
	t := &FetchTask{asset: a, dstDir: dstDir}
	zap.S().Infof("Processing the asset %s, remaining/total: %d/%d", a, len(upscalableAssets), len(assets))

	err = api.GetVariants(t.asset, *t.dstDir)
	if err != nil {
		zap.S().Warnw("Failed to fetch variants of asset",
			"asset", t.asset,
			"error", err,
		)
		return nil
	}
	zap.S().Infow("Fetched variants of asset",
		"asset", t.asset,
	)
	return t
}

func (t *FetchTask) ToUpscaleTask() *UpscaleTask {
	imgV := t.asset.Variants["img"]
	alphaV, hasAlpha := t.asset.Variants["alpha"]
	dstVs := make([]*dto.Variant, 0)
	for _, v := range t.asset.UpscalableVariants {
		v := v
		dstVs = append(dstVs, &v)
	}
	if hasAlpha {
		return &UpscaleTask{
			colorVariant: &imgV,
			alphaVariant: &alphaV,
			dstVariant:   dstVs,
			dstDir:       t.dstDir,
		}
	} else {
		return &UpscaleTask{
			colorVariant: &imgV,
			dstVariant:   dstVs,
			dstDir:       t.dstDir,
		}
	}
}
