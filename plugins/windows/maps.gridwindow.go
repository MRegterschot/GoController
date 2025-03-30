package windows

import (
	"slices"

	"github.com/MRegterschot/GbxRemoteGo/structs"
	"github.com/MRegterschot/GoController/app"
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

func (mgw *MapsGridWindow) HandleRemoveAnswer(login string, data any, _ any) {
	c := app.GetGoController()
	mapInfo := data.(structs.TMMapInfo)

	if err := c.Server.Client.RemoveMap(mapInfo.FileName); err != nil {
		go c.ChatError("Error removing map", err, login)
		return
	}

	// Remove item from grid
	for i, item := range mgw.Items {
		if item.(structs.TMMapInfo).FileName == mapInfo.FileName {
			mgw.Items = slices.Delete(mgw.Items, i, i+1)
			break
		}
	}

	// Remove actions from actions
	for key := range mgw.Actions {
		if key == "remove_"+mapInfo.UId {
			delete(mgw.Actions, key)
			break
		} else if key == "queue_"+mapInfo.UId {
			delete(mgw.Actions, key)
			break
		}
	}

	mgw.Refresh()
	go c.Chat("#Primary#Map removed", login)
}

func (mgw *MapsGridWindow) HandleQueueAnswer(login string, data any, _ any) {
	mapInfo := data.(structs.TMMapInfo)

	app.GetCommandManager().ExecuteCommand(login, "/queue", []string{mapInfo.FileName}, false)
}