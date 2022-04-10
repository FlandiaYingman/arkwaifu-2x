package app

import (
	"github.com/flandiayingman/arkwaifu-2x/internal/app/api"
	"github.com/flandiayingman/arkwaifu-2x/internal/app/dto"
	"github.com/pterm/pterm"
)

func FetchAssets() ([]dto.Asset, error) {
	assets, err := fetchAssets()
	if err != nil {
		return nil, err
	}
	if len(assets) == 0 {
		pterm.Info.Println("No upscalable assets found. Exiting...")
		return nil, nil
	}

	_, err = fetchVariants(assets)
	if err != nil {
		return nil, err
	}

	return assets, nil
}

func fetchAssets() ([]dto.Asset, error) {
	spinner, _ := spinner.Start("Fetching the assets metadata from the server...")

	assets, err := api.GetAssetsUpscalable()
	upscalable, err := api.GetAssetsUpscalable()
	if err != nil {
		_ = spinner.Stop()
		pterm.Error.Println("Failed to fetch the assets metadata from the server.")
		return nil, err
	}

	_ = spinner.Stop()
	pterm.Success.Println("Successfully fetched the assets metadata.")
	pterm.Info.Printfln("Total assets: %d. Upscalable assets: %d", len(assets), len(upscalable))
	return upscalable, nil
}
func fetchVariants(assets []dto.Asset) ([]dto.Variant, error) {
	progressbar, _ := progressbar.
		WithTitle("Fetching the assets file from the server...").
		WithTotal(len(assets)).
		Start()

	vChan, errChan, err := api.GetVariants(assets, dir)
	if err != nil {
		return nil, err
	}

	vs := make([]dto.Variant, 0, len(assets))
	for v := range vChan {
		vs = append(vs, v)
		progressbar.Increment()
	}
	if err, _ := <-errChan; err != nil {
		_, _ = progressbar.Stop()
		pterm.Error.Printfln("Failed to fetch the assets files.")
		return nil, err
	}

	pterm.Success.Println("Successfully fetched the assets files.")
	return vs, nil
}
