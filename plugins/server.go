package plugins

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/MRegterschot/GoController/app"
	"github.com/MRegterschot/GoController/models"
)

type ServerPlugin struct {
	app.BasePlugin
	Name         string
	Dependencies []string
	Loaded       bool
}

func CreateServerPlugin() *ServerPlugin {
	return &ServerPlugin{
		Name:         "Server",
		Dependencies: []string{},
		Loaded:       false,
		BasePlugin:   app.GetBasePlugin(),
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
	if len(args) < 1 {
		err := p.GoController.Server.Client.SetServerPassword("")
		if err != nil {
			go p.GoController.Chat("Error removing server password: "+err.Error(), login)
			return
		}
		go p.GoController.Chat("Server password removed", login)
		return
	}

	err := p.GoController.Server.Client.SetServerPassword(args[0])
	if err != nil {
		go p.GoController.Chat("Error setting server password: "+err.Error(), login)
		return
	}
	go p.GoController.Chat("Server password set", login)
}

func (p *ServerPlugin) getPasswordCommand(login string, args []string) {
	password, err := p.GoController.Server.Client.GetServerPassword()
	if err != nil {
		go p.GoController.Chat("Error getting server password: "+err.Error(), login)
		return
	}
	go p.GoController.Chat("Server password: "+password, login)
}

func (p *ServerPlugin) setPasswordSpectatorCommand(login string, args []string) {
	if len(args) < 1 {
		err := p.GoController.Server.Client.SetServerPasswordForSpectator("")
		if err != nil {
			go p.GoController.Chat("Error removing spectator password: "+err.Error(), login)
			return
		}
		go p.GoController.Chat("Spectator password removed", login)
		return
	}

	err := p.GoController.Server.Client.SetServerPasswordForSpectator(args[0])
	if err != nil {
		go p.GoController.Chat("Error setting spectator password: "+err.Error(), login)
		return
	}
	go p.GoController.Chat("Spectator password set", login)
}

func (p *ServerPlugin) getPasswordSpectatorCommand(login string, args []string) {
	password, err := p.GoController.Server.Client.GetServerPasswordForSpectator()
	if err != nil {
		go p.GoController.Chat("Error getting spectator password: "+err.Error(), login)
		return
	}
	go p.GoController.Chat("Spectator password: "+password, login)
}

func (p *ServerPlugin) nameCommand(login string, args []string) {
	if len(args) < 1 {
		name, err := p.GoController.Server.Client.GetServerName()
		if err != nil {
			go p.GoController.Chat("Error getting server name: "+err.Error(), login)
			return
		}
		go p.GoController.Chat("Server name: "+name, login)
		return
	}

	name := strings.Join(args, " ")

	err := p.GoController.Server.Client.SetServerName(name)
	if err != nil {
		go p.GoController.Chat("Error setting server name: "+err.Error(), login)
		return
	}
	go p.GoController.Chat("Set server name to " + name, login)
}

func (p *ServerPlugin) commentCommand(login string, args []string) {
	if len(args) < 1 {
		comment, err := p.GoController.Server.Client.GetServerComment()
		if err != nil {
			go p.GoController.Chat("Error getting server comment: "+err.Error(), login)
			return
		}
		go p.GoController.Chat("Server comment: "+comment, login)
		return
	}

	comment := strings.Join(args, " ")

	err := p.GoController.Server.Client.SetServerComment(comment)
	if err != nil {
		go p.GoController.Chat("Error setting server comment: "+err.Error(), login)
		return
	}
	go p.GoController.Chat("Set server comment to " + comment, login)
}

func (p *ServerPlugin) hideCommand(login string, args []string) {
	if len(args) < 1 {
		hidden, err := p.GoController.Server.Client.GetHideServer()
		if err != nil {
			go p.GoController.Chat("Error getting server hidden status: "+err.Error(), login)
			return
		}
		if hidden == 1 || hidden == 2 {
			go p.GoController.Chat("Server is hidden", login)
		} else {
			go p.GoController.Chat("Server is not hidden", login)
		}
		return
	}

	hidden := 0
	if args[0] == "1" || args[0] == "true" {
		hidden = 1
	}

	err := p.GoController.Server.Client.SetHideServer(hidden)
	if err != nil {
		go p.GoController.Chat("Error setting server hidden status: "+err.Error(), login)
		return
	}
	if hidden == 1 {
		go p.GoController.Chat("Server is now hidden", login)
	} else {
		go p.GoController.Chat("Server is no longer hidden", login)
	}
}

func (p *ServerPlugin) maxPlayersCommand(login string, args []string) {
	if len(args) < 1 {
		maxPlayers, err := p.GoController.Server.Client.GetMaxPlayers()
		if err != nil {
			go p.GoController.Chat("Error getting max players: "+err.Error(), login)
			return
		}
		go p.GoController.Chat(fmt.Sprintf("Current max: %d, Next max: %d", maxPlayers.CurrentValue, maxPlayers.NextValue), login)
		return
	}

	maxPlayers, err := strconv.Atoi(args[0])
	if err != nil {
		go p.GoController.Chat("Invalid max players", login)
		return
	}

	err = p.GoController.Server.Client.SetMaxPlayers(maxPlayers)
	if err != nil {
		go p.GoController.Chat("Error setting max players: "+err.Error(), login)
		return
	}
	go p.GoController.Chat(fmt.Sprintf("Set max players to %d", maxPlayers), login)
}

func (p *ServerPlugin) maxSpectatorsCommand(login string, args []string) {
	if len(args) < 1 {
		maxSpectators, err := p.GoController.Server.Client.GetMaxSpectators()
		if err != nil {
			go p.GoController.Chat("Error getting max spectators: "+err.Error(), login)
			return
		}
		go p.GoController.Chat(fmt.Sprintf("Current max: %d, Next max: %d", maxSpectators.CurrentValue, maxSpectators.NextValue), login)
		return
	}

	maxSpectators, err := strconv.Atoi(args[0])
	if err != nil {
		go p.GoController.Chat("Invalid max spectators", login)
		return
	}

	err = p.GoController.Server.Client.SetMaxSpectators(maxSpectators)
	if err != nil {
		go p.GoController.Chat("Error setting max spectators: "+err.Error(), login)
		return
	}
	go p.GoController.Chat(fmt.Sprintf("Set max spectators to %d", maxSpectators), login)
}

func (p *ServerPlugin) keepPlayerSlotsCommand(login string, args []string) {
	if len(args) < 1 {
		keepPlayerSlots, err := p.GoController.Server.Client.IsKeepingPlayerSlots()
		if err != nil {
			go p.GoController.Chat("Error getting keep player slots: "+err.Error(), login)
			return
		}

		if keepPlayerSlots {
			go p.GoController.Chat("Keep player slots is enabled", login)
		} else {
			go p.GoController.Chat("Keep player slots is disabled", login)
		}
		return
	}

	keepPlayerSlots := args[0] == "1" || args[0] == "true"

	if err := p.GoController.Server.Client.KeepPlayerSlots(keepPlayerSlots); err != nil {
		go p.GoController.Chat("Error setting keep player slots: "+err.Error(), login)
		return
	}

	if keepPlayerSlots {
		go p.GoController.Chat("Keep player slots is enabled", login)
	} else {
		go p.GoController.Chat("Keep player slots is disabled", login)
	}
}

func (p *ServerPlugin) hornsCommand(login string, args []string) {
	if len(args) < 1 {
		disableHorns, err := p.GoController.Server.Client.AreHornsDisabled()
		if err != nil {
			go p.GoController.Chat("Error getting horns status: "+err.Error(), login)
			return
		}

		if disableHorns {
			go p.GoController.Chat("Horns are disabled", login)
		} else {
			go p.GoController.Chat("Horns are enabled", login)
		}
		return
	}

	disableHorns := false
	if args[0] == "1" || args[0] == "true" {
		disableHorns = true
	}

	err := p.GoController.Server.Client.DisableHorns(disableHorns)
	if err != nil {
		go p.GoController.Chat("Error setting horns status: "+err.Error(), login)
		return
	}

	if disableHorns {
		go p.GoController.Chat("Horns are disabled", login)
	} else {
		go p.GoController.Chat("Horns are enabled", login)
	}
}

func (p *ServerPlugin) serviceAnnouncesCommand(login string, args []string) {
	if len(args) < 1 {
		disableServiceAnnounces, err := p.GoController.Server.Client.AreServiceAnnouncesDisabled()
		if err != nil {
			go p.GoController.Chat("Error getting service announces status: "+err.Error(), login)
			return
		}

		if disableServiceAnnounces {
			go p.GoController.Chat("Service announces are disabled", login)
		} else {
			go p.GoController.Chat("Service announces are enabled", login)
		}
		return
	}

	enableServiceAnnounces := false
	if args[0] == "1" || args[0] == "true" {
		enableServiceAnnounces = true
	}

	err := p.GoController.Server.Client.DisableServiceAnnounces(!enableServiceAnnounces)
	if err != nil {
		go p.GoController.Chat("Error setting service announces status: "+err.Error(), login)
		return
	}

	if enableServiceAnnounces {
		go p.GoController.Chat("Service announces are enabled", login)
		} else {
		go p.GoController.Chat("Service announces are disabled", login)
	}
}

func (p *ServerPlugin) skinsCommand(login string, args []string) {
	if len(args) < 1 {
		disableSkins, err := p.GoController.Server.Client.AreProfileSkinsDisabled()
		if err != nil {
			go p.GoController.Chat("Error getting skins status: "+err.Error(), login)
			return
		}

		if disableSkins {
			go p.GoController.Chat("Skins are disabled", login)
		} else {
			go p.GoController.Chat("Skins are enabled", login)
		}
		return
	}

	enableSkins := false
	if args[0] == "1" || args[0] == "true" {
		enableSkins = true
	}

	err := p.GoController.Server.Client.DisableProfileSkins(!enableSkins)
	if err != nil {
		go p.GoController.Chat("Error setting skins status: "+err.Error(), login)
		return
	}

	if enableSkins {
		go p.GoController.Chat("Skins are enabled", login)
	} else {
		go p.GoController.Chat("Skins are disabled", login)
	}
}

func init() {
	serverPlugin := CreateServerPlugin()
	app.GetPluginManager().PreLoadPlugin(serverPlugin)
}