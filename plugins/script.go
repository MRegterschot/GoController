package plugins

import (
	"fmt"
	"strconv"
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
		Name:     "//modesetting",
		Callback: m.modeSettingCommand,
		Admin:    true,
		Help:     "Set mode setting",
	})

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//modesettings",
		Callback: m.modeSettingsCommand,
		Admin:    true,
		Help:     "Get mode settings",
	})

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//loadmatchsettings",
		Callback: m.loadMatchSettingsCommand,
		Admin:    true,
		Help:     "Load match settings",
	})

	return nil
}

func (m *ScriptPlugin) Unload() error {
	return nil
}

func (m *ScriptPlugin) modeSettingCommand(login string, args []string) {
	if len(args) < 2 {
		go m.GoController.Chat("Usage: //modesetting [*setting] [*value]", login)
		return
	}

	setting := args[0]
	valueStr := strings.Join(args[1:], " ")

	var err error
	switch valueStr {
	case "true":
		err = m.GoController.Server.Client.SetModeScriptSettings(map[string]interface{}{
			setting: true,
		})
	case "false":
		err = m.GoController.Server.Client.SetModeScriptSettings(map[string]interface{}{
			setting: false,
		})
	default:
		// Try to convert valueStr to an integer
		val, err := strconv.Atoi(valueStr)
		if err != nil {
			err = m.GoController.Server.Client.SetModeScriptSettings(map[string]interface{}{
				setting: valueStr,
			})
		} else {
			err = m.GoController.Server.Client.SetModeScriptSettings(map[string]interface{}{
				setting: val,
			})
		}
	}

	if err != nil {
		go m.GoController.Chat("Error setting mode settings: "+err.Error(), login)
		return
	} else {
		go m.GoController.Chat("Setting "+setting+" set to "+valueStr, login)
	}
}

func (m *ScriptPlugin) modeSettingsCommand(login string, args []string) {
	settings, err := m.GoController.Server.Client.GetModeScriptSettings()
	if err != nil {
		go m.GoController.Chat("Error getting mode settings: "+err.Error(), login)
		return
	}

	message := "Mode settings:\n"
	for key, value := range settings {
		message += fmt.Sprintf("%s: %v\n", key, value)
	}

	go m.GoController.Chat(message, login)

	window := ui.NewWindow(&login)
	window.Title = "Mode settings"
	
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

func init() {
	scriptPlugin := CreateScriptPlugin()
	app.GetPluginManager().PreLoadPlugin(scriptPlugin)
}
