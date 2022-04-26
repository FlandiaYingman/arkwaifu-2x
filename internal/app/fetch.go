package app

import (
	"math/rand"
	"time"

	"github.com/flandiayingman/arkwaifu-2x/internal/app/api"
	"github.com/flandiayingman/arkwaifu-2x/internal/app/dto"
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
	// Select asset randomly to avoid 'race' to some degree.
	// The 'race' refers to, like, multiple clients are upscaling the same asset,
	// though at last they will get the same result, it's just a waste of time.
	a := upscalableAssets[rand.Intn(len(upscalableAssets))]
	t := &FetchTask{asset: a, dstDir: dstDir}
	zap.S().Infof("Processing the asset %s, remaning/total: %d/%d", a, len(upscalableAssets), len(assets))

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
