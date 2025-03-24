package plugins

import (
	"strconv"
	"time"

	"github.com/MRegterschot/GoController/app"
	"github.com/MRegterschot/GoController/models"
)

type JukeboxPlugin struct {
	Name         string
	Dependencies []string
	Loaded       bool

	Queue []models.QueueMap
}

func CreateJukeboxPlugin() *JukeboxPlugin {
	return &JukeboxPlugin{
		Name:         "Jukebox",
		Dependencies: []string{},
		Loaded:       false,
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
	c := app.GetGoController()

	if len(args) < 1 {
		if mapInfo, err := c.Server.Client.GetNextMapInfo(); err != nil {
			go c.Chat("Error getting next map info", login)
		} else {
			go c.Chat("Next map: "+mapInfo.Name, login)
		}
		return
	}

	index, err := strconv.Atoi(args[0])
	if err != nil {
		go c.Chat("Invalid index", login)
		return
	}

	err = c.Server.Client.SetNextMapIndex(index)
	if err != nil {
		go c.Chat("Error setting next map", login)
		return
	}

	go c.Chat("Next map set to index "+args[0], login)
}

func (p *JukeboxPlugin) previousCommand(login string, args []string) {
	c := app.GetGoController()

	previousMap := c.MapManager.PreviousMap
	if previousMap == nil {
		go c.Chat("No previous map", login)
		return
	}

	if previousMap.UId == c.MapManager.CurrentMap.UId {
		go c.Chat("Previous map is current map", login)
		return
	}

	p.Queue = append([]models.QueueMap{{
		Name:             previousMap.Name,
		UId:              previousMap.UId,
		FileName:         previousMap.FileName,
		Author:           previousMap.Author,
		AuthorNickname:   previousMap.AuthorNickname,
		QueuedBy:         login,
		QueuedByNickname: c.PlayerManager.GetPlayer(login).NickName,
		QueuedAt:         time.Now(),
	}}, p.Queue...)

	go c.Server.Client.NextMap(false)
}

func (p *JukeboxPlugin) jumpCommand(login string, args []string) {
	c := app.GetGoController()

	if len(args) < 1 {
		go c.Chat("Usage: //jump [*index]", login)
		return
	}

	index, err := strconv.Atoi(args[0])
	if err != nil {
		go c.Chat("Invalid index", login)
		return
	}

	err = c.Server.Client.JumpToMapIndex(index)
	if err != nil {
		go c.Chat("Error jumping to map", login)
		return
	}

	go c.Chat("Jumped to map index "+args[0], login)
}

func init() {
	jukeboxPlugin := CreateJukeboxPlugin()
	app.GetPluginManager().PreLoadPlugin(jukeboxPlugin)
}
