package main

import (
	"context"

	"github.com/merynayr/PingerVK/pinger/internal/app"
	"github.com/merynayr/PingerVK/pkg/logger"
)

func main() {
	ctx := context.Background()

	a, err := app.NewApp(ctx)
	if err != nil {
		logger.Error("failed to init app: %s", err.Error())
	}

	err = a.Run()
	if err != nil {
		logger.Error("failed to run app: %s", err.Error())
	}
}
