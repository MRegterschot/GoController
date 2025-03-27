package widgets

import (
	"errors"
	"strings"
	"sync"

	"github.com/MRegterschot/GoController/app"
	"github.com/MRegterschot/GoController/ui"
	"github.com/MRegterschot/GoController/utils"
	"slices"
)

type Action struct {
	Name    string
	Icon    string
	Command string
}

type AdminControlsWidget struct {
	*ui.Widget
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
		"Controls": acw.Controls,
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

func (acw *AdminControlsWidget) RemoveAction(actionName string) error {
	if !utils.Includes(acw.Actions, actionName) {
		return errors.New("Action does not exist")
	}

	delete(acw.Actions, actionName)
	acw.Controls = removeActionByName(acw.Controls, actionName)
	
	acw.reload()

	return nil
}

func (acw *AdminControlsWidget) executeAction(login string, data any, _ any) {
	command := data.(string)
	cmd := strings.Split(command, " ")[0]
	params := strings.Split(command, " ")[1:]

	app.GetCommandManager().ExecuteCommand(login, cmd, params, true)
}

func removeActionByName(actions []Action, name string) []Action {
	for i, action := range actions {
		if action.Name == name {
			return slices.Delete(actions, i, i+1)
		}
	}
	return actions
}

func init() {
	adminControlsWidget := GetAdminControlsWidget()
	app.GetPluginManager().PreLoadPlugin(adminControlsWidget)
}
