package main

import (
	"time"

	"github.com/MRegterschot/GoController/app"
	"github.com/MRegterschot/GoController/config"
	_ "github.com/MRegterschot/GoController/plugins"
	"github.com/MRegterschot/GoController/utils"
)

func main() {
	config.Setup()

	controller := app.GetGoController()
	controller.Start()

	go utils.MemoryChecker(5 * time.Minute)

	select {}
}