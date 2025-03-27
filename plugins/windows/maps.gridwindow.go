package windows

import "github.com/MRegterschot/GoController/ui"

type MapsGridWindow struct {
	*ui.GridWindow
}

func CreateMapsGridWindow(login *string) *MapsGridWindow {
	mgw := &MapsGridWindow{
		GridWindow: ui.NewGridWindow(login),
	}
	
	mgw.SetTemplate("maps/map.jet")
	mgw.Grid.Rows = 4
	
	return mgw
}

