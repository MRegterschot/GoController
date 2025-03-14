package plugins

import (
	"strconv"

	"github.com/MRegterschot/GoController/app"
	"github.com/MRegterschot/GoController/models"
)

type TeamsPlugin struct {
	app.BasePlugin
	Name         string
	Dependencies []string
	Loaded       bool
}

func CreateTeamsPlugin() *TeamsPlugin {
	return &TeamsPlugin{
		Name:         "Teams",
		Dependencies: []string{},
		Loaded:       false,
		BasePlugin:   app.GetBasePlugin(),
	}
}

func (p *TeamsPlugin) Load() error {
	commandManager := app.GetCommandManager()

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//forcedteams",
		Callback: p.forcedTeamsCommand,
		Admin:    true,
		Help:     "Manage forced teams",
	})

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//forceteam",
		Callback: p.forcePlayerTeamCommand,
		Admin:    true,
		Help:     "Force player team",
	})

	return nil
}

func (p *TeamsPlugin) Unload() error {
	commandManager := app.GetCommandManager()

	commandManager.RemoveCommand("//forcedteams")
	commandManager.RemoveCommand("//forceteam")
	
	return nil
}

func (p *TeamsPlugin) forcedTeamsCommand(login string, args []string) {
	if len(args) < 1 {
		if forcedTeams, err := p.GoController.Server.Client.GetForcedTeams(); err != nil {
			go p.GoController.Chat("Error getting forced teams", login)
		} else {
			if forcedTeams {
				go p.GoController.Chat("Forced teams are enabled", login)
			} else {
				go p.GoController.Chat("Forced teams are disabled", login)
			}
		}
		return
	}

	forcedTeams := args[0] == "true" || args[0] == "1"
	if err := p.GoController.Server.Client.SetForcedTeams(forcedTeams); err != nil {
		go p.GoController.Chat("Error setting forced teams", login)
		return
	}

	if forcedTeams {
		go p.GoController.Chat("Forced teams are enabled", login)
	} else {
		go p.GoController.Chat("Forced teams are disabled", login)
	}
}

func (p *TeamsPlugin) forcePlayerTeamCommand(login string, args []string) {
	if len(args) < 2 {
		go p.GoController.Chat("Usage: //forceteam [*login] [*team]", login)
		return
	}

	team, err := strconv.Atoi(args[1])
	if err != nil {
		go p.GoController.Chat("Invalid team", login)
		return
	}

	if err := p.GoController.Server.Client.ForcePlayerTeam(args[0], team); err != nil {
		go p.GoController.Chat("Error forcing player team: "+err.Error(), login)
		return
	}

	go p.GoController.Chat("Player team forced", login)
}

func init() {
	teamsPlugin := CreateTeamsPlugin()
	app.GetPluginManager().PreLoadPlugin(teamsPlugin)
}
