package plugins

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/MRegterschot/GoController/app"
	"github.com/MRegterschot/GoController/models"
)

type ServerPlugin struct {
	Name         string
	Dependencies []string
	Loaded       bool
}

func CreateServerPlugin() *ServerPlugin {
	return &ServerPlugin{
		Name:         "Server",
		Dependencies: []string{},
		Loaded:       false,
	}
}

func (p *ServerPlugin) Load() error {
	commandManager := app.GetCommandManager()

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//setpassword",
		Callback: p.setPasswordCommand,
		Admin:    true,
		Help:     "Set server password",
		Aliases:  []string{"//setpw"},
	})

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//getpassword",
		Callback: p.getPasswordCommand,
		Admin:    true,
		Help:     "Get server password",
		Aliases:  []string{"//getpw"},
	})

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//setpasswordspectator",
		Callback: p.setPasswordSpectatorCommand,
		Admin:    true,
		Help:     "Set spectator password",
		Aliases:  []string{"//setpwspec"},
	})

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//getpasswordspectator",
		Callback: p.getPasswordSpectatorCommand,
		Admin:    true,
		Help:     "Get spectator password",
		Aliases:  []string{"//getpwspec"},
	})

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//name",
		Callback: p.nameCommand,
		Admin:    true,
		Help:     "Manage server name",
	})

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//comment",
		Callback: p.commentCommand,
		Admin:    true,
		Help:     "Manage server comment",
	})

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//hide",
		Callback: p.hideCommand,
		Admin:    true,
		Help:     "Manage server visibility",
	})

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//maxplayers",
		Callback: p.maxPlayersCommand,
		Admin:    true,
		Help:     "Manage max players",
	})

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//maxspectators",
		Callback: p.maxSpectatorsCommand,
		Admin:    true,
		Help:     "Manage max spectators",
	})

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//keepplayerslots",
		Callback: p.keepPlayerSlotsCommand,
		Admin:    true,
		Help:     "Manage keep player slots",
		Aliases: []string{"//keepslots"},
	})

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//horns",
		Callback: p.hornsCommand,
		Admin:    true,
		Help:     "Manage horns",
	})

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//serviceannounces",
		Callback: p.serviceAnnouncesCommand,
		Admin:    true,
		Help:     "Manage service announces",
		Aliases: []string{"//announces"},
	})

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//skins",
		Callback: p.skinsCommand,
		Admin:    true,
		Help:     "Manage profile skins",
	})

	return nil
}

func (p *ServerPlugin) Unload() error {
	commandManager := app.GetCommandManager()

	commandManager.RemoveCommand("//setpassword")
	commandManager.RemoveCommand("//getpassword")
	commandManager.RemoveCommand("//setpasswordspectator")
	commandManager.RemoveCommand("//getpasswordspectator")
	commandManager.RemoveCommand("//name")
	commandManager.RemoveCommand("//comment")
	commandManager.RemoveCommand("//hide")
	commandManager.RemoveCommand("//maxplayers")
	commandManager.RemoveCommand("//maxspectators")
	commandManager.RemoveCommand("//keepplayerslots")
	commandManager.RemoveCommand("//horns")
	commandManager.RemoveCommand("//serviceannounces")
	commandManager.RemoveCommand("//skins")

	return nil
}

func (p *ServerPlugin) setPasswordCommand(login string, args []string) {
	c := app.GetGoController()
	
	if len(args) < 1 {
		err := c.Server.Client.SetServerPassword("")
		if err != nil {
			go c.Chat("Error removing server password: "+err.Error(), login)
			return
		}
		go c.Chat("Server password removed", login)
		return
	}

	err := c.Server.Client.SetServerPassword(args[0])
	if err != nil {
		go c.Chat("Error setting server password: "+err.Error(), login)
		return
	}
	go c.Chat("Server password set", login)
}

func (p *ServerPlugin) getPasswordCommand(login string, args []string) {
	c := app.GetGoController()
	
	password, err := c.Server.Client.GetServerPassword()
	if err != nil {
		go c.Chat("Error getting server password: "+err.Error(), login)
		return
	}
	go c.Chat("Server password: "+password, login)
}

func (p *ServerPlugin) setPasswordSpectatorCommand(login string, args []string) {
	c := app.GetGoController()
	
	if len(args) < 1 {
		err := c.Server.Client.SetServerPasswordForSpectator("")
		if err != nil {
			go c.Chat("Error removing spectator password: "+err.Error(), login)
			return
		}
		go c.Chat("Spectator password removed", login)
		return
	}

	err := c.Server.Client.SetServerPasswordForSpectator(args[0])
	if err != nil {
		go c.Chat("Error setting spectator password: "+err.Error(), login)
		return
	}
	go c.Chat("Spectator password set", login)
}

func (p *ServerPlugin) getPasswordSpectatorCommand(login string, args []string) {
	c := app.GetGoController()
	
	password, err := c.Server.Client.GetServerPasswordForSpectator()
	if err != nil {
		go c.Chat("Error getting spectator password: "+err.Error(), login)
		return
	}
	go c.Chat("Spectator password: "+password, login)
}

func (p *ServerPlugin) nameCommand(login string, args []string) {
	c := app.GetGoController()
	
	if len(args) < 1 {
		name, err := c.Server.Client.GetServerName()
		if err != nil {
			go c.Chat("Error getting server name: "+err.Error(), login)
			return
		}
		go c.Chat("Server name: "+name, login)
		return
	}

	name := strings.Join(args, " ")

	err := c.Server.Client.SetServerName(name)
	if err != nil {
		go c.Chat("Error setting server name: "+err.Error(), login)
		return
	}
	go c.Chat("Set server name to " + name, login)
}

func (p *ServerPlugin) commentCommand(login string, args []string) {
	c := app.GetGoController()
	
	if len(args) < 1 {
		comment, err := c.Server.Client.GetServerComment()
		if err != nil {
			go c.Chat("Error getting server comment: "+err.Error(), login)
			return
		}
		go c.Chat("Server comment: "+comment, login)
		return
	}

	comment := strings.Join(args, " ")

	err := c.Server.Client.SetServerComment(comment)
	if err != nil {
		go c.Chat("Error setting server comment: "+err.Error(), login)
		return
	}
	go c.Chat("Set server comment to " + comment, login)
}

func (p *ServerPlugin) hideCommand(login string, args []string) {
	c := app.GetGoController()
	
	if len(args) < 1 {
		hidden, err := c.Server.Client.GetHideServer()
		if err != nil {
			go c.Chat("Error getting server hidden status: "+err.Error(), login)
			return
		}
		if hidden == 1 || hidden == 2 {
			go c.Chat("Server is hidden", login)
		} else {
			go c.Chat("Server is not hidden", login)
		}
		return
	}

	hidden := 0
	if args[0] == "1" || args[0] == "true" {
		hidden = 1
	}

	err := c.Server.Client.SetHideServer(hidden)
	if err != nil {
		go c.Chat("Error setting server hidden status: "+err.Error(), login)
		return
	}
	if hidden == 1 {
		go c.Chat("Server is now hidden", login)
	} else {
		go c.Chat("Server is no longer hidden", login)
	}
}

func (p *ServerPlugin) maxPlayersCommand(login string, args []string) {
	c := app.GetGoController()
	
	if len(args) < 1 {
		maxPlayers, err := c.Server.Client.GetMaxPlayers()
		if err != nil {
			go c.Chat("Error getting max players: "+err.Error(), login)
			return
		}
		go c.Chat(fmt.Sprintf("Current max: %d, Next max: %d", maxPlayers.CurrentValue, maxPlayers.NextValue), login)
		return
	}

	maxPlayers, err := strconv.Atoi(args[0])
	if err != nil {
		go c.Chat("Invalid max players", login)
		return
	}

	err = c.Server.Client.SetMaxPlayers(maxPlayers)
	if err != nil {
		go c.Chat("Error setting max players: "+err.Error(), login)
		return
	}
	go c.Chat(fmt.Sprintf("Set max players to %d", maxPlayers), login)
}

func (p *ServerPlugin) maxSpectatorsCommand(login string, args []string) {
	c := app.GetGoController()
	
	if len(args) < 1 {
		maxSpectators, err := c.Server.Client.GetMaxSpectators()
		if err != nil {
			go c.Chat("Error getting max spectators: "+err.Error(), login)
			return
		}
		go c.Chat(fmt.Sprintf("Current max: %d, Next max: %d", maxSpectators.CurrentValue, maxSpectators.NextValue), login)
		return
	}

	maxSpectators, err := strconv.Atoi(args[0])
	if err != nil {
		go c.Chat("Invalid max spectators", login)
		return
	}

	err = c.Server.Client.SetMaxSpectators(maxSpectators)
	if err != nil {
		go c.Chat("Error setting max spectators: "+err.Error(), login)
		return
	}
	go c.Chat(fmt.Sprintf("Set max spectators to %d", maxSpectators), login)
}

func (p *ServerPlugin) keepPlayerSlotsCommand(login string, args []string) {
	c := app.GetGoController()
	
	if len(args) < 1 {
		keepPlayerSlots, err := c.Server.Client.IsKeepingPlayerSlots()
		if err != nil {
			go c.Chat("Error getting keep player slots: "+err.Error(), login)
			return
		}

		if keepPlayerSlots {
			go c.Chat("Keep player slots is enabled", login)
		} else {
			go c.Chat("Keep player slots is disabled", login)
		}
		return
	}

	keepPlayerSlots := args[0] == "1" || args[0] == "true"

	if err := c.Server.Client.KeepPlayerSlots(keepPlayerSlots); err != nil {
		go c.Chat("Error setting keep player slots: "+err.Error(), login)
		return
	}

	if keepPlayerSlots {
		go c.Chat("Keep player slots is enabled", login)
	} else {
		go c.Chat("Keep player slots is disabled", login)
	}
}

func (p *ServerPlugin) hornsCommand(login string, args []string) {
	c := app.GetGoController()
	
	if len(args) < 1 {
		disableHorns, err := c.Server.Client.AreHornsDisabled()
		if err != nil {
			go c.Chat("Error getting horns status: "+err.Error(), login)
			return
		}

		if disableHorns {
			go c.Chat("Horns are disabled", login)
		} else {
			go c.Chat("Horns are enabled", login)
		}
		return
	}

	disableHorns := false
	if args[0] == "1" || args[0] == "true" {
		disableHorns = true
	}

	err := c.Server.Client.DisableHorns(disableHorns)
	if err != nil {
		go c.Chat("Error setting horns status: "+err.Error(), login)
		return
	}

	if disableHorns {
		go c.Chat("Horns are disabled", login)
	} else {
		go c.Chat("Horns are enabled", login)
	}
}

func (p *ServerPlugin) serviceAnnouncesCommand(login string, args []string) {
	c := app.GetGoController()
	
	if len(args) < 1 {
		disableServiceAnnounces, err := c.Server.Client.AreServiceAnnouncesDisabled()
		if err != nil {
			go c.Chat("Error getting service announces status: "+err.Error(), login)
			return
		}

		if disableServiceAnnounces {
			go c.Chat("Service announces are disabled", login)
		} else {
			go c.Chat("Service announces are enabled", login)
		}
		return
	}

	enableServiceAnnounces := false
	if args[0] == "1" || args[0] == "true" {
		enableServiceAnnounces = true
	}

	err := c.Server.Client.DisableServiceAnnounces(!enableServiceAnnounces)
	if err != nil {
		go c.Chat("Error setting service announces status: "+err.Error(), login)
		return
	}

	if enableServiceAnnounces {
		go c.Chat("Service announces are enabled", login)
		} else {
		go c.Chat("Service announces are disabled", login)
	}
}

func (p *ServerPlugin) skinsCommand(login string, args []string) {
	c := app.GetGoController()
	
	if len(args) < 1 {
		disableSkins, err := c.Server.Client.AreProfileSkinsDisabled()
		if err != nil {
			go c.Chat("Error getting skins status: "+err.Error(), login)
			return
		}

		if disableSkins {
			go c.Chat("Skins are disabled", login)
		} else {
			go c.Chat("Skins are enabled", login)
		}
		return
	}

	enableSkins := false
	if args[0] == "1" || args[0] == "true" {
		enableSkins = true
	}

	err := c.Server.Client.DisableProfileSkins(!enableSkins)
	if err != nil {
		go c.Chat("Error setting skins status: "+err.Error(), login)
		return
	}

	if enableSkins {
		go c.Chat("Skins are enabled", login)
	} else {
		go c.Chat("Skins are disabled", login)
	}
}

func init() {
	serverPlugin := CreateServerPlugin()
	app.GetPluginManager().PreLoadPlugin(serverPlugin)
}