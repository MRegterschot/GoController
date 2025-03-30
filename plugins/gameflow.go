package plugins

import (
	"fmt"

	"github.com/MRegterschot/GoController/app"
	"github.com/MRegterschot/GoController/models"
	"github.com/MRegterschot/GoController/plugins/widgets"
)

type GameFlowPlugin struct {
	Name         string
	Dependencies []string
	Loaded       bool
}

func CreateGameFlowPlugin() *GameFlowPlugin {
	return &GameFlowPlugin{
		Name:         "GameFlow",
		Dependencies: []string{},
		Loaded:       false,
	}
}

func (p *GameFlowPlugin) Load() error {
	commandManager := app.GetCommandManager()

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//skip",
		Callback: p.skipCommand,
		Admin:    true,
		Help:     "Skips map",
	})

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//restart",
		Callback: p.restartCommand,
		Admin:    true,
		Help:     "Restarts map",
		Aliases:  []string{"//res"},
	})

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//mode",
		Callback: p.modeCommand,
		Admin:    true,
		Help:     "Get or set gamemode",
	})

	acw := widgets.GetAdminControlsWidget()

	acw.AddAction(widgets.Action{
		Name:    "Skip",
		Icon:    "Skip",
		Command: "//skip",
	})

	acw.AddAction(widgets.Action{
		Name:    "Restart",
		Icon:    "Restart",
		Command: "//restart",
	})

	return nil
}

func (p *GameFlowPlugin) Unload() error {
	commandManager := app.GetCommandManager()

	commandManager.RemoveCommand("//skip")
	commandManager.RemoveCommand("//restart")
	commandManager.RemoveCommand("//mode")

	acw := widgets.GetAdminControlsWidget()

	acw.RemoveAction("Skip")
	acw.RemoveAction("Restart")

	return nil
}

func (p *GameFlowPlugin) skipCommand(login string, args []string) {
	app.GetGoController().Server.Client.NextMap(true)
}

func (p *GameFlowPlugin) restartCommand(login string, args []string) {
	dontClearCupScores := false
	if len(args) > 0 && args[0] == "true" {
		dontClearCupScores = true
	}

	app.GetGoController().Server.Client.RestartMap(dontClearCupScores)
}

func (p *GameFlowPlugin) modeCommand(login string, args []string) {
	c := app.GetGoController()

	if len(args) < 1 {
		if mode, err := c.Server.Client.GetScriptName(); err != nil {
			go c.Chat("#Error#Error getting mode, "+err.Error(), login)
		} else {
			go c.Chat(fmt.Sprintf("#Primary#Current mode: #White#%s, #Primary#Next mode: #White#%s", mode.CurrentValue, mode.NextValue), login)
		}
		return
	}

	if err := c.Server.Client.SetScriptName(args[0]); err != nil {
		go c.Chat("#Error#Error setting mode, "+err.Error(), login)
		return
	}

	go c.Chat("#Primary#Mode set to #White#"+args[0], login)
}

func init() {
	gameFlowPlugin := CreateGameFlowPlugin()
	app.GetPluginManager().PreLoadPlugin(gameFlowPlugin)
}
