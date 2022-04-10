package app

import (
	"github.com/pterm/pterm"
)

const (
	dir = "./arkwaifu-2x-tmp"
)

var spinner = pterm.DefaultSpinner.WithRemoveWhenDone(false)
var progressbar = pterm.DefaultProgressbar.WithRemoveWhenDone(false)
