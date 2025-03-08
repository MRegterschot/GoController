package plugins

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

	slw.AddApplyButtons()
	slw.UpdateItems = slw.updateItems

	return slw
}

func (slw *ScriptListWindow) AddApplyButtons() {
	slw.Actions["apply"] = app.GetUIManager().AddAction(slw.onApply, nil)
	slw.Actions["cancel"] = app.GetUIManager().AddAction(slw.Destroy, nil)
}

func (slw *ScriptListWindow) onApply(login string, _ interface{}, entries interface{}) {
	slw.updateItems(slw.Items, entries)

	var items = make(map[string]interface{})
	for _, item := range slw.Items {
		items[item.Name] = utils.ConvertStringToType(item.Value)
	}

	err := app.GetClient().SetModeScriptSettings(items)
	if err != nil {
		go app.GetGoController().Chat("Error setting mode settings: "+err.Error(), login)
		return
	}

	go app.GetUIManager().DestroyManialink(slw.Manialink)
	go app.GetGoController().Chat("Mode settings applied", login)
}

func (slw *ScriptListWindow) updateItems(items []ui.ListItem, updatedItems interface{}) {
	updatedList, ok := updatedItems.([]structs.TMSEntryVal)
	if !ok {
		return
	}

	for _, updatedItem := range updatedList {
		for i, item := range items {
			if item.Name == updatedItem.Name {
				items[i].Value = updatedItem.Value
				break
			}
		}
	}
}