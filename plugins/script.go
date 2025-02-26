package plugins

import (
	"strings"

	"github.com/MRegterschot/GoController/app"
	"github.com/MRegterschot/GoController/models"
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
		Callback: m.ModeSettingCommand,
		Admin:    true,
		Help:     "Set mode settings",
	})

	return nil
}

func (m *ScriptPlugin) Unload() error {
	return nil
}

func (m *ScriptPlugin) ModeSettingCommand(login string, args []string) {
	if len(args) < 2 {
		go m.GoController.Chat("Usage: //modesetting [*setting] [*value]", login)
		return
	}

	setting := args[0]
	value := strings.Join(args[1:], " ")

	err := m.GoController.Server.Client.SetModeScriptSettings(map[string]interface{}{
		setting: value,
	})

	if err != nil {
		go m.GoController.Chat("Error setting mode settings: "+err.Error(), login)
		return
	} else {
		go m.GoController.Chat("Setting "+setting+" set to "+value, login)
	}
}

func init() {
	scriptPlugin := CreateScriptPlugin()
	app.GetPluginManager().PreLoadPlugin(scriptPlugin)
}
