package plugins

import (
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

func (m *ServerPlugin) Load() error {
	commandManager := app.GetCommandManager()

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//setpassword",
		Callback: m.setPasswordCommand,
		Admin:    true,
		Help:     "Set server password",
		Aliases:  []string{"//setpw"},
	})

	return nil
}

func (m *ServerPlugin) Unload() error {
	commandManager := app.GetCommandManager()

	commandManager.RemoveCommand("//setpassword")

	return nil
}

func (m *ServerPlugin) setPasswordCommand(login string, args []string) {
	if len(args) < 1 {
		err := m.GoController.Server.Client.SetServerPassword("")
		if err != nil {
			go m.GoController.Chat("Error removing server password: "+err.Error(), login)
			return
		}
		go m.GoController.Chat("Server password removed", login)
		return
	}

	err := m.GoController.Server.Client.SetServerPassword(args[0])
	if err != nil {
		go m.GoController.Chat("Error setting server password: "+err.Error(), login)
		return
	}
	go m.GoController.Chat("Server password set", login)
}

func init() {
	serverPlugin := CreateServerPlugin()
	app.GetPluginManager().PreLoadPlugin(serverPlugin)
}