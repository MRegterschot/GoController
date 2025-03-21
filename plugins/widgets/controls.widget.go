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

type ControlsWidget struct {
	*ui.Widget
	app.BasePlugin
	Name         string
	Dependencies []string
	Loaded       bool
	Actions      []Action
}

var (
	cwInstance *ControlsWidget
	cwOnce     sync.Once
)

func GetControlsWidget() *ControlsWidget {
	cwOnce.Do(func() {
		widget := ui.NewWidget("controls/controls.jet")

		widget.Pos = app.UIPos{
			X: -160,
			Y: -39,
		}

		cwInstance = &ControlsWidget{
			Name:         "ControlsWidget",
			Dependencies: []string{},
			Loaded:       false,
			BasePlugin:   app.GetBasePlugin(),
			Widget:       widget,
		}
	})

	return cwInstance
}

func (cw *ControlsWidget) Load() error {
	cw.Display()

	return nil
}

func (cw *ControlsWidget) Unload() error {
	cw.Destroy()

	return nil
}

func (cw *ControlsWidget) AddAction(action Action) error {
	if utils.Includes(cw.Actions, action) {
		return errors.New("Action already exists")
	}

	cw.Actions = append(cw.Actions, action)

	return nil
}

func (cw *ControlsWidget) RemoveAction(action Action) error {
	if !utils.Includes(cw.Actions, action) {
		return errors.New("Action does not exist")
	}

	for i, a := range cw.Actions {
		if a == action {
			cw.Actions = slices.Delete(cw.Actions, i, i+1)
			break
		}
	}

	return nil
}

func init() {
	controlsWidget := GetControlsWidget()
	app.GetPluginManager().PreLoadPlugin(controlsWidget)
}
