package app

import (
	"fmt"

	"github.com/flandiayingman/arkwaifu-2x/internal/app/api"
	"github.com/flandiayingman/arkwaifu-2x/internal/app/dto"
	"github.com/pterm/pterm"
)

func SubmitAssets(variants []dto.Variant) error {
	progressbar, _ := progressbar.
		WithTitle("Submitting the upscaled variants...").
		WithTotal(len(variants)).
		Start()

	vChan, errChan, err := api.PostVariants(variants, dir)
	if err != nil {
		return err
	}
	for v := range vChan {
		progressbar.Increment()
		pterm.Success.Println(fmt.Sprintf("Successfully submitted the upscaled asset: %s.", v.String()))
	}
	for err := range errChan {
		pterm.Error.Println("Failed to submit the upscaled asset.")
		return err
	}
	return nil
}
