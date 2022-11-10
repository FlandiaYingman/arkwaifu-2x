package api

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/cavaliergopher/grab/v3"
	"github.com/flandiayingman/arkwaifu-2x/internal/app/dto"
	"github.com/flandiayingman/arkwaifu-2x/internal/app/util/fileutil"
	"github.com/flandiayingman/arkwaifu-2x/internal/app/verify"
	"github.com/go-resty/resty/v2"
	. "github.com/samber/lo"
)

var rClient = resty.New()
var gClient = grab.NewClient()

var apiUrl = os.Getenv("API_URL")
var assetsUrl = fmt.Sprintf("%s/asset/assets", apiUrl)

func init() {
	rClient.SetRetryCount(3)
}

func GetAssets() ([]dto.Asset, error) {
	resp, err := rClient.
		R().
		SetHeader("Cache-Control", "no-cache").
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
func GetAssetsUpscalable() ([]dto.Asset, []dto.Asset, error) {
	assets, err := GetAssets()
	if err != nil {
		return nil, nil, err
	}

	upscalable := Filter(assets, func(a dto.Asset, _ int) bool { return len(a.UpscalableVariants) > 0 })
	return assets, upscalable, nil
}

func getFile(url string, dst string) error {
	done, _ := verify.Verify(dst)
	if done {
		return nil
	}

	resp, err := rClient.
		R().
		Get(url)
	if err != nil {
		return fmt.Errorf("failed to get file from %s: %w", url, err)
	}

	err = fileutil.MkFileFromBytes(dst, resp.Body())
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", dst, err)
	}

	err = verify.Done(dst)
	if err != nil {
		return fmt.Errorf("cannot verify %s: %w", dst, err)
	}

	return nil
}

func GetVariants(assets dto.Asset, dstDir string) error {
	vimg := assets.Variants["img"]
	valpha, hasAlpha := assets.Variants["alpha"]
	err := getFile(vimg.FileURL(apiUrl), vimg.ActualPath(dstDir))
	if err != nil {
		return err
	}
	if hasAlpha {
		err = getFile(valpha.FileURL(apiUrl), valpha.ActualPath(dstDir))
		if err != nil {
			return err
		}
	}
	return nil
}
func PostVariants(variants []*dto.Variant, srcDir string) error {
	for _, v := range variants {
		url := v.URL(apiUrl)
		resp, err := rClient.
			R().
			SetFile("file", v.ActualPath(srcDir)).
			Post(url)
		if err != nil {
			return fmt.Errorf("failed to post variant %s: %w", v, err)
		}
		if !resp.IsSuccess() {
			return fmt.Errorf("failed to post variant %s: response %d %q", v, resp.StatusCode(), resp.String())
		}
	}
	return nil
}
