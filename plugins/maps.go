package plugins

import (
	"github.com/MRegterschot/GoController/app"
	"github.com/MRegterschot/GoController/models"
	"github.com/MRegterschot/GoController/plugins/windows"
	"github.com/MRegterschot/GoController/utils"
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

	isAdmin := utils.Includes(*c.Admins, login)
	window.IsAdmin = &isAdmin

	for _, m := range c.MapManager.Maps {
		if isAdmin {
			window.Actions["remove_"+m.UId] = app.GetUIManager().AddAction(window.HandleRemoveAnswer, m)
		}
		window.Actions["queue_"+m.UId] = app.GetUIManager().AddAction(window.HandleQueueAnswer, m)
		window.Items = append(window.Items, m)
	}

	go window.Display()
}

func init() {
	mapsPlugin := CreateMapsPlugin()
	app.GetPluginManager().PreLoadPlugin(mapsPlugin)
}
