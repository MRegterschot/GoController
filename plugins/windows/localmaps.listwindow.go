package windows

import (
	"github.com/MRegterschot/GoController/ui"
)

type LocalMapsListWindow struct {
	*ui.ListWindow
}

func CreateLocalMapsListWindow(login *string) *LocalMapsListWindow {
	lmlw := &LocalMapsListWindow{
		ListWindow: ui.NewListWindow(login),
	}

	return lmlw
}