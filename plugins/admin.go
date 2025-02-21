package plugins

import (
	"github.com/MRegterschot/GoController/app"
)

type AdminPlugin struct {
	app.BasePlugin
	Name string
	Dependencies []string
	Loaded bool
}

func CreateAdminPlugin() *AdminPlugin {
	return &AdminPlugin{
		Name: "Admin",
		Dependencies: []string{},
		Loaded: false,
		BasePlugin: app.BasePlugin{
			CommandManager: app.GetCommandManager(),
			SettingsManager: app.GetSettingsManager(),
			GoController: app.GetGoController(),
		},
	}
}

func (m *AdminPlugin) Load() error {
	// commandManager := app.GetCommandManager()

	return nil
}

func (m *AdminPlugin) Unload() error {
	return nil
}

func init() {
	mapsPlugin := CreateAdminPlugin()
	app.GetPluginManager().PreLoadPlugin(mapsPlugin)
}