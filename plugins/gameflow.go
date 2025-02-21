package plugins

import (
	"fmt"

	"github.com/MRegterschot/GoController/app"
	"github.com/MRegterschot/GoController/models"
)

type GameFlowPlugin struct {
	app.BasePlugin
	Name string
	Dependencies []string
	Loaded bool
}

func CreateGameFlowPlugin() *GameFlowPlugin {
	return &GameFlowPlugin{
		Name: "GameFlow",
		Dependencies: []string{},
		Loaded: false,
		BasePlugin: app.BasePlugin{
			CommandManager: app.GetCommandManager(),
			SettingsManager: app.GetSettingsManager(),
			GoController: app.GetGoController(),
		},
	}
}

func (m *GameFlowPlugin) Load() error {
	commandManager := app.GetCommandManager()

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//skip",
		Callback: m.SkipCommand,
		Admin:    true,
		Help:     "Skips map",
	})

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//restart",
		Callback: m.RestartCommand,
		Admin:    true,
		Help:     "Restarts map",
	})

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//mode",
		Callback: m.ModeCommand,
		Admin:    true,
		Help:     "Get or set gamemode",
	})

	return nil
}

func (m *GameFlowPlugin) Unload() error {
	return nil
}

func (m *GameFlowPlugin) SkipCommand(login string, args []string) {
	m.GoController.Server.Client.NextMap()
}

func (m *GameFlowPlugin) RestartCommand(login string, args []string) {
	m.GoController.Server.Client.RestartMap()
}

func (m *GameFlowPlugin) ModeCommand(login string, args []string) {
	if len(args) < 1 {
		if mode, err := m.GoController.Server.Client.GetScriptName(); err != nil {
			go m.GoController.Chat("Error getting mode: " + err.Error(), login)
			return
		} else {
			go m.GoController.Chat(fmt.Sprintf("Current mode: %v, Next mode: %v", mode.CurrentValue, mode.NextValue), login)
			return
		}
	}

	if err := m.GoController.Server.Client.SetScriptName(args[0]); err != nil {
		go m.GoController.Chat("Error setting mode: " + err.Error(), login)
		return
	}

	go m.GoController.Chat("Mode set to " + args[0], login)
}

func init() {
	gameFlowPlugin := CreateGameFlowPlugin()
	app.GetPluginManager().PreLoadPlugin(gameFlowPlugin)
}