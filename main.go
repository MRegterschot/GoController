package main

import (
	"time"

	"github.com/MRegterschot/GoController/app"
	"github.com/MRegterschot/GoController/config"
	_ "github.com/MRegterschot/GoController/plugins"
	_ "github.com/MRegterschot/GoController/plugins/widgets"
	"github.com/MRegterschot/GoController/utils"
	"go.uber.org/zap"
)

func main() {
	config.Setup()
	defer zap.L().Sync()
	utils.SetTheme()

	controller := app.GetGoController()
	controller.Start()

	go utils.MemoryChecker(5 * time.Minute)

	select {}
}