package widgets

import (
	"errors"
	"sync"

	"slices"

	"github.com/MRegterschot/GoController/app"
	"github.com/MRegterschot/GoController/ui"
	"github.com/MRegterschot/GoController/utils"
)

type Action struct {
	Name    string
	Icon    string
	Command string
}

type AdminControlsWidget struct {
	*ui.Widget
	app.BasePlugin
	Name         string
	Dependencies []string
	Loaded       bool
	Controls     []Action
}

var (
	acwInstance *AdminControlsWidget
	acwOnce     sync.Once
)

func GetAdminControlsWidget() *AdminControlsWidget {
	acwOnce.Do(func() {
		widget := ui.NewWidget("controls/admincontrols.jet")

		widget.Pos = app.UIPos{
			X: -160,
			Y: -39,
		}

		acwInstance = &AdminControlsWidget{
			Name:         "AdminControlsWidget",
			Dependencies: []string{},
			Loaded:       false,
			BasePlugin:   app.GetBasePlugin(),
			Widget:       widget,
		}
	})

	return acwInstance
}

func (acw *AdminControlsWidget) Load() error {
	acw.reload()

	return nil
}

func (acw *AdminControlsWidget) Unload() error {
	acw.Destroy()

	return nil
}

func (acw *AdminControlsWidget) reload() {
	acw.Data = map[string]any{
		"Controls": acw.Actions,
	}

	acw.Display()
}

func (acw *AdminControlsWidget) AddAction(action Action) error {
	if utils.Includes(acw.Actions, action) {
		return errors.New("Action already exists")
	}

	acw.Actions[action.Name] = app.GetUIManager().AddAction(acw.executeAction, action.Command)
	acw.Controls = append(acw.Controls, action)
	acw.reload()

	return nil
}

func (acw *AdminControlsWidget) RemoveAction(action Action) error {
	if !utils.Includes(acw.Controls, action) {
		return errors.New("Action does not exist")
	}

	for i, a := range acw.Controls {
		if a == action {
			acw.Controls = slices.Delete(acw.Controls, i, i+1)
			break
		}
	}

	acw.reload()

	return nil
}

func (acw *AdminControlsWidget) executeAction(login string, data any, _ any) {
	command := data.(string)

	app.GetCommandManager().ExecuteCommand(login, command, nil, true)
}

func init() {
	adminControlsWidget := GetAdminControlsWidget()
	app.GetPluginManager().PreLoadPlugin(adminControlsWidget)
}
