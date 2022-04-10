package api

import (
	"github.com/flandiayingman/arkwaifu-2x/internal/app/dto"
	"github.com/flandiayingman/arkwaifu-2x/internal/app/sr"
	. "github.com/samber/lo"
)

type RawAsset struct {
	Kind     string       `json:"kind"`
	Name     string       `json:"name"`
	Variants []RawVariant `json:"variants"`
}
type RawVariant struct {
	Variant  string `json:"variant"`
	Filename string `json:"filename"`
}

func (a RawAsset) toAsset() dto.Asset {
	asset := dto.Asset{
		Kind: a.Kind,
		Name: a.Name,
	}
	asset.Variants = KeyBy(Map(a.Variants,
		func(v RawVariant, _ int) dto.Variant {
			return v.toVariant(a)
		}),
		func(v dto.Variant) string {
			return v.Variant
		})
	asset.UpscalableVariants = sr.UpscalableVariants(Values(asset.Variants))
	return asset
}
func (v RawVariant) toVariant(a RawAsset) dto.Variant {
	variant := dto.Variant{
		Kind:     a.Kind,
		Name:     a.Name,
		Variant:  v.Variant,
		Filename: v.Filename,
	}
	return variant
}
