package windows

import (
	"github.com/MRegterschot/GbxRemoteGo/structs"
	"github.com/MRegterschot/GoController/app"
	"github.com/MRegterschot/GoController/ui"
	"github.com/MRegterschot/GoController/utils"
)

type ScriptListWindow struct {
	*ui.ListWindow
}

func CreateScriptListWindow(login *string) *ScriptListWindow {
	slw := &ScriptListWindow{
		ListWindow: ui.NewListWindow(login),
	}

	slw.AddActionButtons()
	slw.UpdateItems = slw.updateItems

	return slw
}

func (slw *ScriptListWindow) AddActionButtons() {
	slw.Actions["apply"] = app.GetUIManager().AddAction(slw.onApply, nil)
	slw.Actions["cancel"] = app.GetUIManager().AddAction(slw.Destroy, nil)
}

func (slw *ScriptListWindow) onApply(login string, _ any, entries any) {
	slw.updateItems(slw.Items, entries)

	var items = make(map[string]any)
	for _, item := range slw.Items {
		if key, ok := item[0].(string); ok {
			items[key] = utils.ConvertStringToType(item[2].(string))
		}
	}

	err := app.GetClient().SetModeScriptSettings(items)
	if err != nil {
		go app.GetGoController().Chat("#Error#Error setting mode settings, "+err.Error(), login)
		return
	}

	go app.GetGoController().Chat("#Primary#Mode settings applied", login)
}

func (slw *ScriptListWindow) updateItems(items [][]any, updatedItems any) {
	updatedList, ok := updatedItems.([]structs.TMSEntryVal)
	if !ok {
		return
	}

	for _, updatedItem := range updatedList {
		for i, item := range items {
			if key, ok := item[0].(string); ok && key == updatedItem.Name {
				items[i][2] = updatedItem.Value
				break
			}
		}
	}
}
