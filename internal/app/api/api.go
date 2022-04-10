package api

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/cavaliergopher/grab/v3"
	"github.com/flandiayingman/arkwaifu-2x/internal/app/dto"
	"github.com/go-resty/resty/v2"
	. "github.com/samber/lo"
)

var rClient = resty.New()
var gClient = grab.NewClient()

// TODO: Get it from the environment variables.
const apiUrl = "http://localhost:7080/api/v0"

var (
	assetsUrl = fmt.Sprintf("%s/asset/assets", apiUrl)
)

func GetAssets() ([]dto.Asset, error) {
	resp, err := rClient.
		R().
		Get(assetsUrl)
	if err != nil {
		return nil, err
	}

	var rawAssets []RawAsset
	err = json.Unmarshal(resp.Body(), &rawAssets)
	if err != nil {
		return nil, err
	}

	assets := Map(rawAssets, func(r RawAsset, _ int) dto.Asset { return r.toAsset() })
	return assets, nil
}
func GetAssetsUpscalable() ([]dto.Asset, error) {
	assets, err := GetAssets()
	if err != nil {
		return nil, err
	}

	upscalable := Filter(assets, func(a dto.Asset, _ int) bool { return len(a.UpscalableVariants) > 0 })
	return upscalable, nil
}

func GetVariants(assets []dto.Asset, dstDir string) (<-chan dto.Variant, <-chan error, error) {
	requests := Map(assets, func(a dto.Asset, _ int) *grab.Request {
		v := a.Variants["img"]
		url := v.FileURL(apiUrl)
		dst := filepath.Join(dstDir, v.Path())

		req, err := grab.NewRequest(dst, url)
		if err != nil {
			panic(err)
		}

		req.Tag = v
		return req
	})
	respChan := gClient.DoBatch(8, requests...)

	vChan := make(chan dto.Variant, len(assets))
	errChan := make(chan error)
	go func() {
		defer close(vChan)
		defer close(errChan)
		for resp := range respChan {
			err := resp.Err()
			if err != nil {
				errChan <- err
				return
			}
			v := resp.Request.Tag.(dto.Variant)
			vChan <- v
		}
	}()

	return vChan, errChan, nil
}
func PostVariants(variants []dto.Variant, dir string) (<-chan dto.Variant, <-chan error, error) {
	lock := make(chan struct{}, 8)
	vChan := make(chan dto.Variant, len(variants))
	errChan := make(chan error)

	go func() {
		defer close(vChan)
		defer close(errChan)
		for _, v := range variants {
			lock <- struct{}{}
			go func(v dto.Variant) {
				url := fmt.Sprintf(v.URL(apiUrl))
				resp, err := rClient.
					R().
					SetFile("file", filepath.Join(dir, v.Path())).
					Post(url)

				if err != nil {
					errChan <- err
					return
				}

				if !resp.IsSuccess() {
					errChan <- fmt.Errorf("response %s not success", resp.String())
					return
				}

				vChan <- v
				<-lock
			}(v)
		}
	}()
	return vChan, errChan, nil
}
