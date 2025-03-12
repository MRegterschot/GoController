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
		BasePlugin: app.GetBasePlugin(),
	}
}

func (m *PlayersPlugin) Load() error {
	commandManager := app.GetCommandManager()

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//ban",
		Callback: m.banCommand,
		Admin:    true,
		Help:     "Bans a player",
	})

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//unban",
		Callback: m.unBanCommand,
		Admin:    true,
		Help:     "Unbans a player",
	})

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//banlist",
		Callback: m.banListCommand,
		Admin:    true,
		Help:     "Lists all bans",
	})

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//cleanbanlist",
		Callback: m.cleanBanListCommand,
		Admin:    true,
		Help:     "Cleans the ban list",
		Aliases: []string{"//cbl"},
	})

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//kick",
		Callback: m.kickCommand,
		Admin:    true,
		Help:     "Kicks a player",
	})

	return nil
}

func (m *PlayersPlugin) Unload() error {
	commandManager := app.GetCommandManager()

	commandManager.RemoveCommand("//ban")
	commandManager.RemoveCommand("//unban")
	commandManager.RemoveCommand("//banlist")
	commandManager.RemoveCommand("//cleanbanlist")
	commandManager.RemoveCommand("//kick")

	return nil
}

func (m *PlayersPlugin) banCommand(login string, args []string) {
	if len(args) < 1 {
		go m.GoController.Chat("Usage: //ban [*login] [reason]", login)
		return
	}

	targetLogin := args[0]
	reason := ""
	if len(args) > 1 {
		reason = strings.Join(args[1:], " ")
	}
	if err := m.GoController.Server.Client.Ban(targetLogin, reason); err != nil {
		go m.GoController.Chat("Error banning player: "+err.Error(), login)
		return
	}

	go m.GoController.Chat(fmt.Sprintf("Banned: %s, Reason: %s", targetLogin, reason))
}

func (m *PlayersPlugin) unBanCommand(login string, args []string) {
	if len(args) < 1 {
		go m.GoController.Chat("Usage: //unban [*login]", login)
		return
	}

	targetLogin := args[0]
	if err := m.GoController.Server.Client.UnBan(targetLogin); err != nil {
		go m.GoController.Chat("Error unbanning player: "+err.Error(), login)
		return
	}
	go m.GoController.Chat(fmt.Sprintf("Unbanned: %s", targetLogin))
}

func (m *PlayersPlugin) banListCommand(login string, args []string) {
	banList, err := m.GoController.Server.Client.GetBanList(100, 0)
	if err != nil {
		go m.GoController.Chat("Error getting ban list: "+err.Error(), login)
		return
	}

	if len(banList) == 0 {
		go m.GoController.Chat("No bans found", login)
		return
	}

	logins := make([]string, len(banList))
	for i, ban := range banList {
		logins[i] = ban.Login
	}

	msg := fmt.Sprintf("Bans (%d): %s", len(banList), strings.Join(logins, ", "))

	go m.GoController.Chat(msg, login)
}

func (m *PlayersPlugin) cleanBanListCommand(login string, args []string) {
	if err := m.GoController.Server.Client.CleanBanList(); err != nil {
		go m.GoController.Chat("Error cleaning ban list: "+err.Error(), login)
		return
	}

	go m.GoController.Chat("Ban list cleaned", login)
}

func (m *PlayersPlugin) kickCommand(login string, args []string) {
	if len(args) < 1 {
		go m.GoController.Chat("Usage: //kick [*login] [reason]", login)
		return
	}

	targetLogin := args[0]
	reason := ""
	if len(args) > 1 {
		reason = strings.Join(args[1:], " ")
	}
	m.GoController.Server.Client.Kick(targetLogin, reason)
	go m.GoController.Chat(fmt.Sprintf("Kicked: %s, Reason: %s", targetLogin, reason))
}

func init() {
	playersPlugin := CreatePlayersPlugin()
	app.GetPluginManager().PreLoadPlugin(playersPlugin)
}