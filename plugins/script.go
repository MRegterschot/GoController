package plugins

import (
	"fmt"
	"sort"
	"strings"

	"github.com/MRegterschot/GoController/app"
	"github.com/MRegterschot/GoController/models"
	"github.com/MRegterschot/GoController/plugins/windows"
	"github.com/MRegterschot/GoController/ui"
)

type ScriptPlugin struct {
	Name         string
	Dependencies []string
	Loaded       bool
}

func CreateScriptPlugin() *ScriptPlugin {
	return &ScriptPlugin{
		Name:         "Script",
		Dependencies: []string{},
		Loaded:       false,
	}
}

func (p *ScriptPlugin) Load() error {
	commandManager := app.GetCommandManager()

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//modesettings",
		Callback: p.modeSettingsCommand,
		Admin:    true,
		Help:     "Get mode settings",
		Aliases:  []string{"//ms"},
	})

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//loadmatchsettings",
		Callback: p.loadMatchSettingsCommand,
		Admin:    true,
		Help:     "Load match settings",
		Aliases:  []string{"//lms"},
	})

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//savematchsettings",
		Callback: p.saveMatchSettingsCommand,
		Admin:    true,
		Help:     "Save match settings",
		Aliases:  []string{"//sms"},
	})

	return nil
}

func (p *ScriptPlugin) Unload() error {
	commandManager := app.GetCommandManager()

	commandManager.RemoveCommand("//modesettings")
	commandManager.RemoveCommand("//loadmatchsettings")
	commandManager.RemoveCommand("//savematchsettings")

	return nil
}

func (p *ScriptPlugin) modeSettingsCommand(login string, args []string) {
	c := app.GetGoController()
	
	settings, err := c.Server.Client.GetModeScriptSettings()
	if err != nil {
		go c.ChatError("Error getting mode settings", err, login)
		return
	}

	info, err := c.Server.Client.GetModeScriptInfo()
	if err != nil {
		go c.ChatError("#Error getting mode script info", err, login)
		return
	}

	// Extract and sort keys
	keys := make([]string, 0, len(settings))
	for key := range settings {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	items := make([][]any, 0, len(keys))
	for _, key := range keys {
		var desc string
		for _, item := range info.ParamDescs {
			if item.Name == key && item.Desc != "<hidden>" {
				desc = item.Desc
				break
			}
		}

		items = append(items, []any{
			key,
			desc,
			fmt.Sprintf("%v", settings[key]),
		})
	}

	columns := []ui.Column{
		{Name: "Name", Width: 30},
		{Name: "Description", Width: 40},
		{Name: "Value", Width: 30, Type: "input"},
	}

	window := windows.CreateScriptListWindow(&login)
	window.Title = "Mode settings"
	window.Columns = columns
	window.Items = items

	go window.Display()
}

func (p *ScriptPlugin) loadMatchSettingsCommand(login string, args []string) {
	c := app.GetGoController()
	
	if len(args) < 1 {
		go c.ChatUsage("//loadmatchsettings [*filename]", login)
		return
	}

	filename := args[0]
	_, err := c.Server.Client.LoadMatchSettings("MatchSettings/" + filename)
	if err != nil {
		go c.ChatError("Error loading match settings", err, login)
		return
	}

	go c.Chat("#Primary#Match settings loaded", login)
}

func (p *ScriptPlugin) saveMatchSettingsCommand(login string, args []string) {
	c := app.GetGoController()
	
	file := "tracklist.txt"
	if len(args) > 0 {
		cleanFile, _ := strings.CutSuffix(args[0], ".txt")
		file = cleanFile + ".txt"
	}

	_, err := c.Server.Client.SaveMatchSettings("MatchSettings/" + file)
	if err != nil {
		go c.ChatError("Error saving match settings", err, login)
		return
	}

	go c.Chat("#Primary#Match settings saved to #White#"+file, login)
}

func init() {
	scriptPlugin := CreateScriptPlugin()
	app.GetPluginManager().PreLoadPlugin(scriptPlugin)
}
