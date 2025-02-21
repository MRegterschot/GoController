package plugins

import (
	"fmt"
	"strings"

	"github.com/MRegterschot/GoController/app"
	"github.com/MRegterschot/GoController/models"
)

type PlayersPlugin struct {
	app.BasePlugin
	Name string
	Dependencies []string
	Loaded bool
}

func CreatePlayersPlugin() *PlayersPlugin {
	return &PlayersPlugin{
		Name: "Players",
		Dependencies: []string{},
		Loaded: false,
		BasePlugin: app.BasePlugin{
			CommandManager: app.GetCommandManager(),
			SettingsManager: app.GetSettingsManager(),
			GoController: app.GetGoController(),
		},
	}
}

func (m *PlayersPlugin) Load() error {
	commandManager := app.GetCommandManager()

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//ban",
		Callback: m.BanCommand,
		Admin:    true,
		Help:     "Bans a player",
	})

	return nil
}

func (m *PlayersPlugin) Unload() error {
	return nil
}

func (m *PlayersPlugin) BanCommand(login string, args []string) {
	if len(args) < 1 {
		go m.GoController.Chat("Usage: //ban [*login] [reason]", login)
		return
	}

	targetLogin := args[0]
	reason := ""
	if len(args) > 1 {
		reason = strings.Join(args[1:], " ")
	}
	m.GoController.Server.Client.BanAndBlackList(targetLogin, reason, true)
	go m.GoController.Chat(fmt.Sprintf("Banned: %s, Reason: %s", targetLogin, reason))
}

func init() {
	playersPlugin := CreatePlayersPlugin()
	app.GetPluginManager().PreLoadPlugin(playersPlugin)
}