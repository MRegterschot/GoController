package main

import (
	"github.com/MRegterschot/GoController/app"
	"github.com/MRegterschot/GoController/config"
	"github.com/MRegterschot/GoController/utils"
	_ "github.com/MRegterschot/GoController/plugins"
)

func main() {
	config.Setup()

	controller := app.GetGoController()
	controller.Start()

	go utils.MemoryChecker()

	select {}
}