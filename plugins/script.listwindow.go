package plugins

import (
	"fmt"

	"github.com/MRegterschot/GoController/app"
	"github.com/MRegterschot/GoController/ui"
)

type ScriptListWindow struct {
	ui.ListWindow
}

func CreateScriptListWindow(login *string) *ScriptListWindow {
	slw := &ScriptListWindow{
		ListWindow: *ui.NewListWindow(login),
	}

	slw.AddApplyButtons()

	return slw
}

func (slw *ScriptListWindow) AddApplyButtons() {
	slw.Actions["apply"] = app.GetUIManager().AddAction(slw.onApply, "")
	slw.Actions["cancel"] = app.GetUIManager().AddAction(slw.Destroy, "")
}

func (slw *ScriptListWindow) onApply(_ string, _ interface{}, _ interface{}) {
	fmt.Println("glorp")
}