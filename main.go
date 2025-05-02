/*
    GoController - A trackmania server controller
    Copyright (C) 2025 MRegterschot

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"time"

	"github.com/MRegterschot/GoController/app"
	"github.com/MRegterschot/GoController/config"
	_ "github.com/MRegterschot/GoController/plugins"
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