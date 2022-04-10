package main

// This import must be first. Because we wanna logger be initialized first.
import _ "github.com/flandiayingman/arkwaifu-2x/internal/app/log"

import (
	"github.com/flandiayingman/arkwaifu-2x/internal/app"
	"go.uber.org/zap"
)

func main() {
	log := zap.S()
	defer func() { _ = log.Sync() }()

	assets, err := app.FetchAssets()
	if err != nil {
		log.Panic(err)
		return
	}
	if len(assets) == 0 {
		return
	}

	variants, err := app.UpscaleAssets(assets)
	if err != nil {
		log.Panic(err)
	}

	err = app.SubmitAssets(variants)
	if err != nil {
		log.Panic(err)
	}
}
