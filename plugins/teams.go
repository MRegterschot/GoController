package plugins

import (
	"strconv"

	"github.com/MRegterschot/GoController/app"
	"github.com/MRegterschot/GoController/models"
)

type TeamsPlugin struct {
	Name         string
	Dependencies []string
	Loaded       bool
}

func CreateTeamsPlugin() *TeamsPlugin {
	return &TeamsPlugin{
		Name:         "Teams",
		Dependencies: []string{},
		Loaded:       false,
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
	c := app.GetGoController()
	
	if len(args) < 1 {
		if forcedTeams, err := c.Server.Client.GetForcedTeams(); err != nil {
			go c.ChatError("Error getting forced teams", err, login)
		} else {
			if forcedTeams {
				go c.Chat("#Primary#Forced teams are enabled", login)
			} else {
				go c.Chat("#Primary#Forced teams are disabled", login)
			}
		}
		return
	}

	forcedTeams := args[0] == "true" || args[0] == "1"
	if err := c.Server.Client.SetForcedTeams(forcedTeams); err != nil {
		go c.ChatError("Error setting forced teams", err, login)
		return
	}

	if forcedTeams {
		go c.Chat("#Primary#Forced teams are enabled", login)
	} else {
		go c.Chat("#Primary#Forced teams are disabled", login)
	}
}

func (p *TeamsPlugin) forcePlayerTeamCommand(login string, args []string) {
	c := app.GetGoController()
	
	if len(args) < 2 {
		go c.ChatUsage("//forceteam [*login] [*team]", login)
		return
	}

	team, err := strconv.Atoi(args[1])
	if err != nil {
		go c.ChatError("Invalid team", nil, login)
		return
	}

	if err := c.Server.Client.ForcePlayerTeam(args[0], team); err != nil {
		go c.ChatError("Error forcing player team", err, login)
		return
	}

	go c.Chat("#Primary#Player team forced", login)
}

func init() {
	teamsPlugin := CreateTeamsPlugin()
	app.GetPluginManager().PreLoadPlugin(teamsPlugin)
}
