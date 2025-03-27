package windows

import (
	"github.com/MRegterschot/GoController/ui"
)

type MapsGridWindow struct {
	*ui.GridWindow
	IsAdmin *bool
}

func CreateMapsGridWindow(login *string) *MapsGridWindow {
	mgw := &MapsGridWindow{
		GridWindow: ui.NewGridWindow(login),
	}

	mgw.GridWindow.AddData = mgw.CheckAdmin

	mgw.SetTemplate("maps/map.jet")
	mgw.Grid.Rows = 4
	
	return mgw
}

func (mgw *MapsGridWindow) CheckAdmin() {
	if dataMap, ok := mgw.GridWindow.Data.(map[string]any); ok {
		dataMap["IsAdmin"] = *mgw.IsAdmin
	}
}
