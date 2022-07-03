package main

// This import must be first. Because we wanna logger be initialized first.
import (
	"time"

	"github.com/flandiayingman/arkwaifu-2x/internal/app"
	_ "github.com/flandiayingman/arkwaifu-2x/internal/app/log"
	"github.com/panjf2000/ants/v2"
)

import (
	"go.uber.org/zap"
)

func main() {
	log := zap.S()
	defer func() { _ = log.Sync() }()

	err := Run()
	if err != nil {
		log.Error(err)
	}
}

func Run() (err error) {
	defer func() {
		if err != nil {
			return
		}
		if err0 := recover(); err0 != nil {
			if _, ok := err0.(error); ok {
				err = err0.(error)
			}
		}
	}()

	inPool, err := ants.NewPool(1)
	if err != nil {
		return err
	}
	procPool, err := ants.NewPool(2)
	if err != nil {
		return err
	}
	outPool, err := ants.NewPool(4)
	if err != nil {
		return err
	}

	var task func()
	task = func() {
		fetchTask := app.CreateTask(&app.Dir)
		if fetchTask == nil {
			time.Sleep(15 * time.Second)
			return
		}
		upscaleTask := fetchTask.ToUpscaleTask()
		err := procPool.Submit(func() {
			upscaleTask.Run()
			submitTask := upscaleTask.ToSubmitTask()
			err := outPool.Submit(func() {
				submitTask.Run()
			})
			panic(err)
		})
		if err != nil {
			panic(err)
		}
	}
	for {
		err = inPool.Submit(task)
		if err != nil {
			panic(err)
		}
	}
}
