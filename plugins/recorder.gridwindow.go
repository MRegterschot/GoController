package plugins

import (
	"github.com/MRegterschot/GoController/ui"
)

type RecorderGridWindow struct {
	*ui.GridWindow
}

func CreateRecorderGridWindow(login *string) *RecorderGridWindow {
	rgw := &RecorderGridWindow{
		GridWindow: ui.NewGridWindow(login),
	}

	return rgw
}