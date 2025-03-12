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

func (p *ServerPlugin) Load() error {
	commandManager := app.GetCommandManager()

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//setpassword",
		Callback: p.setPasswordCommand,
		Admin:    true,
		Help:     "Set server password",
		Aliases:  []string{"//setpw"},
	})

	return nil
}

func (p *ServerPlugin) Unload() error {
	commandManager := app.GetCommandManager()

	commandManager.RemoveCommand("//setpassword")

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

func init() {
	serverPlugin := CreateServerPlugin()
	app.GetPluginManager().PreLoadPlugin(serverPlugin)
}