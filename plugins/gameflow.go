package plugins

import (
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

func init() {
	gameFlowPlugin := CreateGameFlowPlugin()
	app.GetPluginManager().PreLoadPlugin(gameFlowPlugin)
}