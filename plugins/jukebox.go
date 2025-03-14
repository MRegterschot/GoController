package plugins

import (
	"strconv"

	"github.com/MRegterschot/GoController/app"
	"github.com/MRegterschot/GoController/models"
)

type JukeboxPlugin struct {
	app.BasePlugin
	Name         string
	Dependencies []string
	Loaded       bool
}

func CreateJukeboxPlugin() *JukeboxPlugin {
	return &JukeboxPlugin{
		Name:         "Jukebox",
		Dependencies: []string{},
		Loaded:       false,
		BasePlugin:   app.GetBasePlugin(),
	}
}

func (p *JukeboxPlugin) Load() error {
	commandManager := app.GetCommandManager()

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//next",
		Callback: p.nextCommand,
		Admin:    true,
		Help:     "Manage next map",
	})

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//jump",
		Callback: p.jumpCommand,
		Admin:    true,
		Help:     "Jump to map",
	})

	return nil
}

func (p *JukeboxPlugin) Unload() error {
	commandManager := app.GetCommandManager()

	commandManager.RemoveCommand("//next")
	commandManager.RemoveCommand("//jump")

	return nil
}

func (p *JukeboxPlugin) nextCommand(login string, args []string) {
	if len(args) < 1 {
		if index, err := p.GoController.Server.Client.GetNextMapIndex(); err != nil {
			go p.GoController.Chat("Error getting next map index", login)
		} else {
			go p.GoController.Chat("Next map index is " + strconv.Itoa(index), login)
		}
		return
	}

	index, err := strconv.Atoi(args[0])
	if err != nil {
		go p.GoController.Chat("Invalid index", login)
		return
	}

	err = p.GoController.Server.Client.SetNextMapIndex(index)
	if err != nil {
		go p.GoController.Chat("Error setting next map", login)
		return
	}

	go p.GoController.Chat("Next map set to index " + args[0], login)
}

func (p *JukeboxPlugin) jumpCommand(login string, args []string) {
	if len(args) < 1 {
		go p.GoController.Chat("Usage: //jump [*index]", login)
		return
	}

	index, err := strconv.Atoi(args[0])
	if err != nil {
		go p.GoController.Chat("Invalid index", login)
		return
	}

	err = p.GoController.Server.Client.JumpToMapIndex(index)
	if err != nil {
		go p.GoController.Chat("Error jumping to map", login)
		return
	}

	go p.GoController.Chat("Jumped to map index " + args[0], login)
}

func init() {
	jukeboxPlugin := CreateJukeboxPlugin()
	app.GetPluginManager().PreLoadPlugin(jukeboxPlugin)
}