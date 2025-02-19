package main

import (
	"github.com/MRegterschot/GoController/config"
)


func main() {
	config.Setup()

	controller := GetController()
	controller.Start()


	select {}
}