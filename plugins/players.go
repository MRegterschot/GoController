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

func (p *PlayersPlugin) Load() error {
	commandManager := app.GetCommandManager()

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//ban",
		Callback: p.banCommand,
		Admin:    true,
		Help:     "Bans a player",
	})

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//unban",
		Callback: p.unBanCommand,
		Admin:    true,
		Help:     "Unbans a player",
	})

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//banlist",
		Callback: p.banListCommand,
		Admin:    true,
		Help:     "Lists all bans",
	})

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//cleanbanlist",
		Callback: p.cleanBanListCommand,
		Admin:    true,
		Help:     "Cleans the ban list",
	})

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//blacklist",
		Callback: p.blackListCommand,
		Admin:    true,
		Help:     "Blacklists a player",
	})

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//unblacklist",
		Callback: p.unBlackListCommand,
		Admin:    true,
		Help:     "Unblacklists a player",
	})

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//loadblacklist",
		Callback: p.loadBlackListCommand,
		Admin:    true,
		Help:     "Loads a black list",
	})

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//saveblacklist",
		Callback: p.saveBlackListCommand,
		Admin:    true,
		Help:     "Saves the black list",
	})

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//cleanblacklist",
		Callback: p.cleanBlackListCommand,
		Admin:    true,
		Help:     "Cleans the black list",
	})

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//kick",
		Callback: p.kickCommand,
		Admin:    true,
		Help:     "Kicks a player",
	})

	return nil
}

func (p *PlayersPlugin) Unload() error {
	commandManager := app.GetCommandManager()

	commandManager.RemoveCommand("//ban")
	commandManager.RemoveCommand("//unban")
	commandManager.RemoveCommand("//banlist")
	commandManager.RemoveCommand("//cleanbanlist")
	commandManager.RemoveCommand("//blacklist")
	commandManager.RemoveCommand("//unblacklist")
	commandManager.RemoveCommand("//loadblacklist")
	commandManager.RemoveCommand("//saveblacklist")
	commandManager.RemoveCommand("//cleanblacklist")
	commandManager.RemoveCommand("//kick")

	return nil
}

func (p *PlayersPlugin) banCommand(login string, args []string) {
	if len(args) < 1 {
		go p.GoController.Chat("Usage: //ban [*login] [reason]", login)
		return
	}

	targetLogin := args[0]
	reason := ""
	if len(args) > 1 {
		reason = strings.Join(args[1:], " ")
	}
	if err := p.GoController.Server.Client.Ban(targetLogin, reason); err != nil {
		go p.GoController.Chat("Error banning player: "+err.Error(), login)
		return
	}

	go p.GoController.Chat(fmt.Sprintf("Banned: %s, Reason: %s", targetLogin, reason))
}

func (p *PlayersPlugin) unBanCommand(login string, args []string) {
	if len(args) < 1 {
		go p.GoController.Chat("Usage: //unban [*login]", login)
		return
	}

	targetLogin := args[0]
	if err := p.GoController.Server.Client.UnBan(targetLogin); err != nil {
		go p.GoController.Chat("Error unbanning player: "+err.Error(), login)
		return
	}
	go p.GoController.Chat(fmt.Sprintf("Unbanned: %s", targetLogin))
}

func (p *PlayersPlugin) banListCommand(login string, args []string) {
	banList, err := p.GoController.Server.Client.GetBanList(100, 0)
	if err != nil {
		go p.GoController.Chat("Error getting ban list: "+err.Error(), login)
		return
	}

	if len(banList) == 0 {
		go p.GoController.Chat("No bans found", login)
		return
	}

	logins := make([]string, len(banList))
	for i, ban := range banList {
		logins[i] = ban.Login
	}

	msg := fmt.Sprintf("Bans (%d): %s", len(banList), strings.Join(logins, ", "))

	go p.GoController.Chat(msg, login)
}

func (p *PlayersPlugin) cleanBanListCommand(login string, args []string) {
	if err := p.GoController.Server.Client.CleanBanList(); err != nil {
		go p.GoController.Chat("Error cleaning ban list: "+err.Error(), login)
		return
	}

	go p.GoController.Chat("Ban list cleaned", login)
}

func (p *PlayersPlugin) blackListCommand(login string, args []string) {
	if len(args) < 1 {
		blackList, err := p.GoController.Server.Client.GetBlackList(100, 0)
		if err != nil {
			go p.GoController.Chat("Error getting black list: "+err.Error(), login)
			return
		}

		if len(blackList) == 0 {
			go p.GoController.Chat("No blacklisted players found", login)
			return
		}

		logins := make([]string, len(blackList))
		for i, black := range blackList {
			logins[i] = black.Login
		}

		msg := fmt.Sprintf("Blacklisted players (%d): %s", len(blackList), strings.Join(logins, ", "))

		go p.GoController.Chat(msg, login)
		return
	}

	targetLogin := args[0]
	if err := p.GoController.Server.Client.BlackList(targetLogin); err != nil {
		go p.GoController.Chat("Error blacklisting player: "+err.Error(), login)
		return
	}

	go p.GoController.Chat(fmt.Sprintf("Blacklisted: %s", targetLogin))
}

func (p *PlayersPlugin) unBlackListCommand(login string, args []string) {
	if len(args) < 1 {
		go p.GoController.Chat("Usage: //unblacklist [*login]", login)
		return
	}

	targetLogin := args[0]
	if err := p.GoController.Server.Client.UnBlackList(targetLogin); err != nil {
		go p.GoController.Chat("Error unblacklisting player: "+err.Error(), login)
		return
	}

	go p.GoController.Chat(fmt.Sprintf("Unblacklisted: %s", targetLogin))
}

func (p *PlayersPlugin) loadBlackListCommand(login string, args []string) {
	if len(args) < 1 {
		go p.GoController.Chat("Usage: //loadblacklist [*file]", login)
		return
	}

	file := args[0]

	if err := p.GoController.Server.Client.LoadBlackList(file); err != nil {
		go p.GoController.Chat("Error loading black list: "+err.Error(), login)
		return
	}

	go p.GoController.Chat(fmt.Sprintf("Black list %s loaded", file), login)
}

func (p *PlayersPlugin) saveBlackListCommand(login string, args []string) {
	if len(args) < 1 {
		go p.GoController.Chat("Usage: //saveblacklist [*file]", login)
		return
	}

	file := args[0]

	if err := p.GoController.Server.Client.SaveBlackList(file); err != nil {
		go p.GoController.Chat("Error saving black list: "+err.Error(), login)
		return
	}

	go p.GoController.Chat(fmt.Sprintf("Black list saved to %s", file), login)
}

func (p *PlayersPlugin) cleanBlackListCommand(login string, args []string) {
	if err := p.GoController.Server.Client.CleanBlackList(); err != nil {
		go p.GoController.Chat("Error cleaning black list: "+err.Error(), login)
		return
	}

	go p.GoController.Chat("Black list cleaned", login)
}

func (p *PlayersPlugin) kickCommand(login string, args []string) {
	if len(args) < 1 {
		go p.GoController.Chat("Usage: //kick [*login] [reason]", login)
		return
	}

	targetLogin := args[0]
	reason := ""
	if len(args) > 1 {
		reason = strings.Join(args[1:], " ")
	}
	p.GoController.Server.Client.Kick(targetLogin, reason)
	go p.GoController.Chat(fmt.Sprintf("Kicked: %s, Reason: %s", targetLogin, reason))
}

func init() {
	playersPlugin := CreatePlayersPlugin()
	app.GetPluginManager().PreLoadPlugin(playersPlugin)
}