package plugins

import (
	"fmt"

	"github.com/MRegterschot/GoController/app"
	"github.com/MRegterschot/GoController/models"
	"github.com/MRegterschot/GoController/plugins/widgets"
)

type GameFlowPlugin struct {
	app.BasePlugin
	Name         string
	Dependencies []string
	Loaded       bool
}

func CreateGameFlowPlugin() *GameFlowPlugin {
	return &GameFlowPlugin{
		Name:         "GameFlow",
		Dependencies: []string{},
		Loaded:       false,
		BasePlugin:   app.GetBasePlugin(),
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

	cw := widgets.GetControlsWidget()

	cw.AddAction(widgets.Action{
		Name:    "Skip",
		Icon:    "Skip",
		Command: "//skip",
	})

	cw.AddAction(widgets.Action{
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

	return nil
}

func (p *GameFlowPlugin) skipCommand(login string, args []string) {
	dontClearCupScores := false
	if len(args) > 0 && args[0] == "true" {
		dontClearCupScores = true
	}
	
	p.GoController.Server.Client.NextMap(dontClearCupScores)
}

func (p *GameFlowPlugin) restartCommand(login string, args []string) {
	dontClearCupScores := false
	if len(args) > 0 && args[0] == "true" {
		dontClearCupScores = true
	}
	
	p.GoController.Server.Client.RestartMap(dontClearCupScores)
}

func (p *GameFlowPlugin) modeCommand(login string, args []string) {
	if len(args) < 1 {
		if mode, err := p.GoController.Server.Client.GetScriptName(); err != nil {
			go p.GoController.Chat("Error getting mode: "+err.Error(), login)
		} else {
			go p.GoController.Chat(fmt.Sprintf("Current mode: %s, Next mode: %s", mode.CurrentValue, mode.NextValue), login)
		}
		return
	}

	if err := p.GoController.Server.Client.SetScriptName(args[0]); err != nil {
		go p.GoController.Chat("Error setting mode: "+err.Error(), login)
		return
	}

	go p.GoController.Chat("Mode set to "+args[0], login)
}

func init() {
	gameFlowPlugin := CreateGameFlowPlugin()
	app.GetPluginManager().PreLoadPlugin(gameFlowPlugin)
}
