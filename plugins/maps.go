package plugins

import (
	"fmt"

	"github.com/MRegterschot/GoController/app"
	"github.com/MRegterschot/GoController/models"
	"github.com/MRegterschot/GoController/plugins/windows"
)

type MapsPlugin struct {
	Name         string
	Dependencies []string
	Loaded       bool
}

func CreateMapsPlugin() *MapsPlugin {
	return &MapsPlugin{
		Name:         "Maps",
		Dependencies: []string{},
		Loaded:       false,
	}
}

func (p *MapsPlugin) Load() error {
	commandManager := app.GetCommandManager()

	commandManager.AddCommand(models.ChatCommand{
		Name:     "/maps",
		Callback: p.mapsCommand,
		Admin:    false,
		Help:     "Shows all available maps",
	})

	return nil
}

func (p *MapsPlugin) Unload() error {
	commandManager := app.GetCommandManager()

	commandManager.RemoveCommand("/maps")

	return nil
}

func (p *MapsPlugin) mapsCommand(login string, args []string) {
	c := app.GetGoController()

	window := windows.CreateMapsGridWindow(&login)
	window.Title = "Maps"
	window.Items = make([]any, 0)

	for _, m := range c.MapManager.Maps {
		window.Actions["remove_"+m.UId] = app.GetUIManager().AddAction(p.handleRemoveAnswer, m)
		window.Actions["queue_"+m.UId] = app.GetUIManager().AddAction(p.handleQueueAnswer, m)
		window.Items = append(window.Items, m)
	}

	go window.Display()
}

func (p *MapsPlugin) handleRemoveAnswer(login string, data any, _ any) {
	fmt.Println(data)
}

func (p *MapsPlugin) handleQueueAnswer(login string, data any, _ any) {
	fmt.Println(data)
}

func init() {
	mapsPlugin := CreateMapsPlugin()
	app.GetPluginManager().PreLoadPlugin(mapsPlugin)
}
