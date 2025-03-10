package plugins

import (
	"fmt"
	"sort"
	"strings"

	"github.com/MRegterschot/GoController/app"
	"github.com/MRegterschot/GoController/models"
	"github.com/MRegterschot/GoController/ui"
)

type ScriptPlugin struct {
	app.BasePlugin
	Name         string
	Dependencies []string
	Loaded       bool
}

func CreateScriptPlugin() *ScriptPlugin {
	return &ScriptPlugin{
		Name:         "Script",
		Dependencies: []string{},
		Loaded:       false,
		BasePlugin:   app.GetBasePlugin(),
	}
}

func (m *ScriptPlugin) Load() error {
	commandManager := app.GetCommandManager()

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//modesettings",
		Callback: m.modeSettingsCommand,
		Admin:    true,
		Help:     "Get mode settings",
		Aliases:  []string{"//ms"},
	})

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//loadmatchsettings",
		Callback: m.loadMatchSettingsCommand,
		Admin:    true,
		Help:     "Load match settings",
		Aliases:  []string{"//lms"},
	})

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//savematchsettings",
		Callback: m.saveMatchSettingsCommand,
		Admin:    true,
		Help:     "Save match settings",
		Aliases:  []string{"//sms"},
	})

	return nil
}

func (m *ScriptPlugin) Unload() error {
	return nil
}

func (m *ScriptPlugin) modeSettingsCommand(login string, args []string) {
	settings, err := m.GoController.Server.Client.GetModeScriptSettings()
	if err != nil {
		go m.GoController.Chat("Error getting mode settings: "+err.Error(), login)
		return
	}

	info, err := m.GoController.Server.Client.GetModeScriptInfo()
	if err != nil {
		go m.GoController.Chat("Error getting mode script info: "+err.Error(), login)
		return
	}

	// Extract and sort keys
	keys := make([]string, 0, len(settings))
	for key := range settings {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	items := make([]ui.ListItem, 0)
	for _, key := range keys {
		var desc string
		for _, item := range info.ParamDescs {
			if item.Name == key && item.Desc != "<hidden>" {
				desc = item.Desc
				break
			}
		}

		items = append(items, ui.ListItem{
			Name:        key,
			Description: desc,
			Value:       fmt.Sprintf("%v", settings[key]),
		})
	}

	window := CreateScriptListWindow(&login)
	window.Title = "Mode settings"
	window.Items = items

	go window.Display()
}

func (m *ScriptPlugin) loadMatchSettingsCommand(login string, args []string) {
	if len(args) < 1 {
		go m.GoController.Chat("Usage: //loadmatchsettings [*filename]", login)
		return
	}

	filename := args[0]
	_, err := m.GoController.Server.Client.LoadMatchSettings("MatchSettings/" + filename)
	if err != nil {
		go m.GoController.Chat("Error loading match settings: "+err.Error(), login)
		return
	}

	go m.GoController.Chat("Match settings loaded", login)
}

func (m *ScriptPlugin) saveMatchSettingsCommand(login string, args []string) {
	file := "tracklist.txt"
	if len(args) > 0 {
		cleanFile, _ := strings.CutSuffix(args[0], ".txt")
		file = cleanFile + ".txt"
	}

	_, err := m.GoController.Server.Client.SaveMatchSettings("MatchSettings/" + file)
	if err != nil {
		go m.GoController.Chat("Error saving match settings: "+err.Error(), login)
		return
	}

	go m.GoController.Chat("Match settings saved", login)
}

func init() {
	scriptPlugin := CreateScriptPlugin()
	app.GetPluginManager().PreLoadPlugin(scriptPlugin)
}
