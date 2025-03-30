package plugins

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/MRegterschot/GoController/app"
	"github.com/MRegterschot/GoController/models"
)

type PlayersPlugin struct {
	Name string
	Dependencies []string
	Loaded bool
}

var statuses = [4]string{"Selectable", "Spectator", "Player", "Spectator but selectable"}

func CreatePlayersPlugin() *PlayersPlugin {
	return &PlayersPlugin{
		Name: "Players",
		Dependencies: []string{},
		Loaded: false,
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
		Help:     "Lists all players on ban list",
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
		Help:     "Adds player to black list",
	})

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//unblacklist",
		Callback: p.unBlackListCommand,
		Admin:    true,
		Help:     "Removes player from black list",
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
		Name:     "//fakeplayer",
		Callback: p.fakePlayerCommand,
		Admin:    true,
		Help:     "Connects or disconnects a fake player",
	})

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//addguest",
		Callback: p.addGuestCommand,
		Admin:    true,
		Help:     "Adds player to guest list",
	})

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//removeguest",
		Callback: p.removeGuestCommand,
		Admin:    true,
		Help:     "Removes player from guest list",
	})

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//guestlist",
		Callback: p.guestListCommand,
		Admin:    true,
		Help:     "Lists all players on guest list",
	})

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//loadguestlist",
		Callback: p.loadGuestListCommand,
		Admin:    true,
		Help:     "Loads a guest list",
	})

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//saveguestlist",
		Callback: p.saveGuestListCommand,
		Admin:    true,
		Help:     "Saves the guest list",
	})

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//cleanguestlist",
		Callback: p.cleanGuestListCommand,
		Admin:    true,
		Help:     "Cleans the guest list",
	})

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//kick",
		Callback: p.kickCommand,
		Admin:    true,
		Help:     "Kicks a player",
	})

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//players",
		Callback: p.getPlayersCommand,
		Admin:    true,
		Help:     "Lists all players",
	})

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//forcestatus",
		Callback: p.forceStatusCommand,
		Admin:    true,
		Help:     "Forces a player to status",
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
	commandManager.RemoveCommand("//fakeplayer")
	commandManager.RemoveCommand("//addguest")
	commandManager.RemoveCommand("//removeguest")
	commandManager.RemoveCommand("//guestlist")
	commandManager.RemoveCommand("//loadguestlist")
	commandManager.RemoveCommand("//saveguestlist")
	commandManager.RemoveCommand("//cleanguestlist")
	commandManager.RemoveCommand("//kick")
	commandManager.RemoveCommand("//players")

	return nil
}

func (p *PlayersPlugin) banCommand(login string, args []string) {
	c := app.GetGoController()

	if len(args) < 1 {
		go c.Chat("#Primary#Usage: #White#//ban [*login] [reason]", login)
		return
	}

	targetLogin := args[0]
	reason := ""
	if len(args) > 1 {
		reason = strings.Join(args[1:], " ")
	}
	if err := c.Server.Client.Ban(targetLogin, reason); err != nil {
		go c.Chat("#Error#Error banning player, "+err.Error(), login)
		return
	}

	go c.Chat(fmt.Sprintf("#Primary#Banned: #White#%s, #Primary#Reason: #White#%s", targetLogin, reason))
}

func (p *PlayersPlugin) unBanCommand(login string, args []string) {
	c := app.GetGoController()

	if len(args) < 1 {
		go c.Chat("#Primary#Usage: #White#//unban [*login]", login)
		return
	}

	targetLogin := args[0]
	if err := c.Server.Client.UnBan(targetLogin); err != nil {
		go c.Chat("#Error#Error unbanning player, "+err.Error(), login)
		return
	}
	go c.Chat(fmt.Sprintf("#Primary#Unbanned #White#%s", targetLogin))
}

func (p *PlayersPlugin) banListCommand(login string, args []string) {
	c := app.GetGoController()

	banList, err := c.Server.Client.GetBanList(100, 0)
	if err != nil {
		go c.Chat("#Error#Error getting ban list, "+err.Error(), login)
		return
	}

	if len(banList) == 0 {
		go c.Chat("#Primary#No bans found", login)
		return
	}

	logins := make([]string, len(banList))
	for i, ban := range banList {
		logins[i] = ban.Login
	}

	msg := fmt.Sprintf("#Primary#Bans (%d): #White#%s", len(banList), strings.Join(logins, ", "))

	go c.Chat(msg, login)
}

func (p *PlayersPlugin) cleanBanListCommand(login string, args []string) {
	c := app.GetGoController()

	if err := c.Server.Client.CleanBanList(); err != nil {
		go c.Chat("#Error#Error cleaning ban list, "+err.Error(), login)
		return
	}

	go c.Chat("#Primary#Ban list cleaned", login)
}

func (p *PlayersPlugin) blackListCommand(login string, args []string) {
	c := app.GetGoController()

	if len(args) < 1 {
		blackList, err := c.Server.Client.GetBlackList(100, 0)
		if err != nil {
			go c.Chat("#Error#Error getting black list, "+err.Error(), login)
			return
		}

		if len(blackList) == 0 {
			go c.Chat("#Primary#No blacklisted players found", login)
			return
		}

		logins := make([]string, len(blackList))
		for i, black := range blackList {
			logins[i] = black.Login
		}

		msg := fmt.Sprintf("#Primary#Blacklisted players (%d): #White#%s", len(blackList), strings.Join(logins, ", "))

		go c.Chat(msg, login)
		return
	}

	targetLogin := args[0]
	if err := c.Server.Client.BlackList(targetLogin); err != nil {
		go c.Chat("#Error#Error blacklisting player, "+err.Error(), login)
		return
	}

	go c.Chat(fmt.Sprintf("#Primary#Blacklisted #White#%s", targetLogin))
}

func (p *PlayersPlugin) unBlackListCommand(login string, args []string) {
	c := app.GetGoController()

	if len(args) < 1 {
		go c.Chat("#Primary#Usage: #White#//unblacklist [*login]", login)
		return
	}

	targetLogin := args[0]
	if err := c.Server.Client.UnBlackList(targetLogin); err != nil {
		go c.Chat("#Error#Error unblacklisting player, "+err.Error(), login)
		return
	}

	go c.Chat(fmt.Sprintf("#Primary#Unblacklisted #White#%s", targetLogin))
}

func (p *PlayersPlugin) loadBlackListCommand(login string, args []string) {
	c := app.GetGoController()

	if len(args) < 1 {
		go c.Chat("#Primary#Usage: #White#//loadblacklist [*file]", login)
		return
	}

	file := args[0]

	if err := c.Server.Client.LoadBlackList(file); err != nil {
		go c.Chat("#Error#Error loading black list, "+err.Error(), login)
		return
	}

	go c.Chat(fmt.Sprintf("#Primary#Black list #White#%s #Primary#loaded", file), login)
}

func (p *PlayersPlugin) saveBlackListCommand(login string, args []string) {
	c := app.GetGoController()

	if len(args) < 1 {
		go c.Chat("#Primary#Usage: #White#//saveblacklist [*file]", login)
		return
	}

	file := args[0]

	if err := c.Server.Client.SaveBlackList(file); err != nil {
		go c.Chat("#Error#Error saving black list, "+err.Error(), login)
		return
	}

	go c.Chat(fmt.Sprintf("#Primary#Black list saved to #White#%s", file), login)
}

func (p *PlayersPlugin) cleanBlackListCommand(login string, args []string) {
	c := app.GetGoController()

	if err := c.Server.Client.CleanBlackList(); err != nil {
		go c.Chat("#Error#Error cleaning black list, "+err.Error(), login)
		return
	}

	go c.Chat("#Primary#Black list cleaned", login)
}

func (p *PlayersPlugin) fakePlayerCommand(login string, args []string) {
	c := app.GetGoController()

	if len(args) > 0 {
		targetLogin := args[0]
		if err := c.Server.Client.DisconnectFakePlayer(targetLogin); err != nil {
			go c.Chat("#Error#Error disconnecting fake player, "+err.Error(), login)
			return
		}

		go c.Chat(fmt.Sprintf("#Primary#Fake player #White#%s #Primary#disconnected", targetLogin), login)
		return
	}

	if err := c.Server.Client.ConnectFakePlayer(); err != nil {
		go c.Chat("#Error#Error connecting fake player, "+err.Error(), login)
		return
	}

	go c.Chat("#Primary#Fake player connected", login)
}

func (p *PlayersPlugin) addGuestCommand(login string, args []string) {
	c := app.GetGoController()
	
	if len(args) < 1 {
		go c.Chat("#Primary#Usage: #White#//addguest [*login]", login)
		return
	}

	targetLogin := args[0]
	if err := c.Server.Client.AddGuest(targetLogin); err != nil {
		go c.Chat("#Error#Error adding guest, "+err.Error(), login)
		return
	}

	go c.Chat(fmt.Sprintf("#Primary#Guest added: #White#%s", targetLogin), login)
}

func (p *PlayersPlugin) removeGuestCommand(login string, args []string) {
	c := app.GetGoController()
	
	if len(args) < 1 {
		go c.Chat("#Primary#Usage: #White#//removeguest [*login]", login)
		return
	}

	targetLogin := args[0]
	if err := c.Server.Client.RemoveGuest(targetLogin); err != nil {
		go c.Chat("#Error#Error removing guest, "+err.Error(), login)
		return
	}

	go c.Chat(fmt.Sprintf("#Primary#Guest removed: #White#%s", targetLogin), login)
}

func (p *PlayersPlugin) guestListCommand(login string, args []string) {
	c := app.GetGoController()
	
	guestList, err := c.Server.Client.GetGuestList(100, 0)
	if err != nil {
		go c.Chat("#Error#Error getting guest list, "+err.Error(), login)
		return
	}

	if len(guestList) == 0 {
		go c.Chat("#Primary#No guests found", login)
		return
	}

	logins := make([]string, len(guestList))
	for i, guest := range guestList {
		logins[i] = guest.Login
	}

	msg := fmt.Sprintf("#Primary#Guests (%d): #White#%s", len(guestList), strings.Join(logins, ", "))

	go c.Chat(msg, login)
}

func (p *PlayersPlugin) loadGuestListCommand(login string, args []string) {
	c := app.GetGoController()
	
	if len(args) < 1 {
		go c.Chat("#Primary#Usage: #White#//loadguestlist [*file]", login)
		return
	}

	file := args[0]

	if err := c.Server.Client.LoadGuestList(file); err != nil {
		go c.Chat("#Error#Error loading guest list, "+err.Error(), login)
		return
	}

	go c.Chat(fmt.Sprintf("#Primary#Guest list #White#%s #Primary#loaded", file), login)
}

func (p *PlayersPlugin) saveGuestListCommand(login string, args []string) {
	c := app.GetGoController()
	
	if len(args) < 1 {
		go c.Chat("#Primary#Usage: #White#//saveguestlist [*file]", login)
		return
	}

	file := args[0]

	if err := c.Server.Client.SaveGuestList(file); err != nil {
		go c.Chat("#Error#Error saving guest list, "+err.Error(), login)
		return
	}

	go c.Chat(fmt.Sprintf("#Primary#Guest list saved to #White#%s", file), login)
}

func (p *PlayersPlugin) cleanGuestListCommand(login string, args []string) {
	c := app.GetGoController()
	
	if err := c.Server.Client.CleanGuestList(); err != nil {
		go c.Chat("#Error#Error cleaning guest list, "+err.Error(), login)
		return
	}

	go c.Chat("#Primary#Guest list cleaned", login)
}

func (p *PlayersPlugin) kickCommand(login string, args []string) {
	c := app.GetGoController()
	
	if len(args) < 1 {
		go c.Chat("#Primary#Usage: #White#//kick [*login] [reason]", login)
		return
	}

	targetLogin := args[0]
	reason := ""
	if len(args) > 1 {
		reason = strings.Join(args[1:], " ")
	}
	c.Server.Client.Kick(targetLogin, reason)
	go c.Chat(fmt.Sprintf("#Primary#Kicked: #White#%s, #Primary#Reason: #White#%s", targetLogin, reason))
}

func (p *PlayersPlugin) getPlayersCommand(login string, args []string) {
	c := app.GetGoController()
	
	players, err := c.Server.Client.GetPlayerList(100, 1)
	if err != nil {
		go c.Chat("#Error#Error getting players, "+err.Error(), login)
		return
	}

	if len(players) == 0 {
		go c.Chat("#Primary#No players found", login)
		return
	}

	nickNames := make([]string, len(players))
	for i, player := range players {
		nickNames[i] = player.NickName
	}

	msg := fmt.Sprintf("#Primary#Players (%d): #White#%s", len(players), strings.Join(nickNames, ", "))

	go c.Chat(msg, login)
}

func (p *PlayersPlugin) forceStatusCommand(login string, args []string) {
	c := app.GetGoController()
	
	if len(args) < 1 {
		go c.Chat("#Primary#Usage: #White#//forcestatus [*login] [status]", login)
		return
	}

	status := 3
	if len(args) > 1 {
		argInt, err := strconv.Atoi(args[1])
		if err != nil || argInt < 0 || argInt > len(statuses) - 1 {
			go c.Chat("#Error#Invalid status", login)
			return
		}

		status = argInt
	}

	targetLogin := args[0]
	if err := c.Server.Client.ForceSpectator(targetLogin, status); err != nil {
		go c.Chat("#Error#Error forcing status, "+err.Error(), login)
		return
	}

	go c.Chat(fmt.Sprintf("#Primary#Forced #White#%s #Primary#into status #White#%s", targetLogin, statuses[status]))
}

func init() {
	playersPlugin := CreatePlayersPlugin()
	app.GetPluginManager().PreLoadPlugin(playersPlugin)
}