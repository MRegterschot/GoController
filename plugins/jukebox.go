package plugins

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/MRegterschot/GbxRemoteGo/gbxclient"
	"github.com/MRegterschot/GbxRemoteGo/structs"
	"github.com/MRegterschot/GoController/app"
	"github.com/MRegterschot/GoController/models"
	"github.com/MRegterschot/GoController/plugins/widgets"
)

type JukeboxPlugin struct {
	Name         string
	Dependencies []string
	Loaded       bool

	Queue []models.QueueMap
}

var (
	jukebox *JukeboxPlugin
	jukeboxOnce sync.Once
)

func CreateJukeboxPlugin() *JukeboxPlugin {
	return &JukeboxPlugin{
		Name:         "Jukebox",
		Dependencies: []string{},
		Loaded:       false,
	}
}

func GetJukeboxPlugin() *JukeboxPlugin {
	jukeboxOnce.Do(func() {
		jukebox = CreateJukeboxPlugin()
	})
	return jukebox
}

func (p *JukeboxPlugin) Load() error {
	commandManager := app.GetCommandManager()

	app.GetClient().ScriptCallbacks["Maniaplanet.Podium_Start"] = append(app.GetClient().ScriptCallbacks["Maniaplanet.Podium_Start"], gbxclient.GbxCallbackStruct[any]{
		Key:  "jukebox",
		Call: p.onEndRace,
	})

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//next",
		Callback: p.nextCommand,
		Admin:    true,
		Help:     "Manage next map",
	})

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//previous",
		Callback: p.previousCommand,
		Admin:    true,
		Help:     "Jump to previous map",
		Aliases:  []string{"//prev"},
	})

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//jump",
		Callback: p.jumpCommand,
		Admin:    true,
		Help:     "Jump to map",
	})

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//requeue",
		Callback: p.requeueCommand,
		Admin:    true,
		Help:     "Requeue current map",
	})

	commandManager.AddCommand(models.ChatCommand{
		Name:     "/queue",
		Callback: p.queueCommand,
		Admin:    false,
		Help:     "Queue map",
	})

	acw := widgets.GetAdminControlsWidget()

	acw.AddAction(widgets.Action{
		Name: "Previous",
		Icon: "Previous",
		Command: "//previous",
	})

	acw.AddAction(widgets.Action{
		Name: "Requeue",
		Icon: "Requeue",
		Command: "//requeue",
	})

	return nil
}

func (p *JukeboxPlugin) Unload() error {
	commandManager := app.GetCommandManager()

	commandManager.RemoveCommand("//next")
	commandManager.RemoveCommand("//previous")
	commandManager.RemoveCommand("//jump")
	commandManager.RemoveCommand("//requeue")
	commandManager.RemoveCommand("/queue")

	acw := widgets.GetAdminControlsWidget()

	acw.RemoveAction("Previous")
	acw.RemoveAction("Requeue")

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

	go c.Chat("Next map set to index "+args[0])
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

	if err := c.Server.Client.ChooseNextMap(previousMap.FileName); err != nil {
		go c.Chat("Error setting previous map", login)
		return
	}

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

	go c.Chat("Jumped to map index "+args[0])
}

func (p *JukeboxPlugin) requeueCommand(login string, args []string) {
	c := app.GetGoController()

	currentMap := c.MapManager.CurrentMap
	if currentMap.UId == "" {
		go c.Chat("No current map", login)
		return
	}

	if len(p.Queue) > 0 && p.Queue[0].UId == currentMap.UId {
		go c.Chat("Map already in queue", login)
		return
	}

	p.QueueMap(currentMap, login)

	go c.Chat("Map requeued")
}

func (p *JukeboxPlugin) queueCommand(login string, args []string) {
	c := app.GetGoController()

	if len(args) < 1 {
		go c.Chat("Usage: /queue [*filename]", login)
		return
	}

	mapInfo, err := c.Server.Client.GetMapInfo(args[0])
	if err != nil {
		go c.Chat("Error getting map info", login)
		return
	}

	if err := p.QueueMap(mapInfo, login); err != nil {
		go c.Chat("Error queuing map", login)
		return
	}

	go c.Chat("Map queued")
}

func (p *JukeboxPlugin) onEndRace(_ any) {
	c := app.GetGoController()

	if len(p.Queue) == 0 {
		mapInfo, err := c.Server.Client.GetNextMapInfo()
		if err != nil {
			return
		}

		go c.Chat(fmt.Sprintf("Next map %s by %s", mapInfo.Name, mapInfo.AuthorNickname))
		return
	}

	nextMap := p.Queue[0]
	p.Queue = p.Queue[1:]

	if err := c.Server.Client.ChooseNextMap(nextMap.FileName); err != nil {
		go c.Chat("Error setting next map")
		return
	}

	go c.Chat(fmt.Sprintf("Next map %s by %s queued by %s", nextMap.Name, nextMap.AuthorNickname, nextMap.QueuedByNickname))
}

func (p *JukeboxPlugin) QueueMap(mapInfo structs.TMMapInfo, login string) error {
	c := app.GetGoController()

	player, err := c.Server.Client.GetPlayerInfo(login)
	if err != nil {
		return err
	}

	var queueMap models.QueueMap
	queueMap.ToQueueMap(mapInfo)
	queueMap.QueuedBy = login
	queueMap.QueuedByNickname = player.NickName

	p.Queue = append(p.Queue, queueMap)

	go c.Chat(fmt.Sprintf("Map %s by %s queued by %s", mapInfo.Name, mapInfo.AuthorNickname, player.NickName))

	return nil
}

func init() {
	jukeboxPlugin := GetJukeboxPlugin()
	app.GetPluginManager().PreLoadPlugin(jukeboxPlugin)
}
