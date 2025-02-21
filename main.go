package main

import (
	"github.com/MRegterschot/GoController/config"
	"github.com/MRegterschot/GoController/utils"
)

func main() {
	config.Setup()

	controller := GetController()
	controller.Start()

	go utils.MemoryChecker()

	select {}
}