package plugins

import (
	"errors"
	"fmt"
	"strings"

	"github.com/MRegterschot/GoController/app"
	"github.com/MRegterschot/GoController/models"
	"github.com/MRegterschot/GoController/ui"
)

type ControllerPlugin struct {
	Name         string
	Dependencies []string
	Loaded       bool
}

func CreateControllerPlugin() *ControllerPlugin {
	return &ControllerPlugin{
		Name:         "Controller",
		Dependencies: []string{},
		Loaded:       false,
	}
}

func (p *ControllerPlugin) Load() error {
	commandManager := app.GetCommandManager()

	commandManager.AddCommand(models.ChatCommand{
		Name:     "/help",
		Callback: p.helpCommand,
		Admin:    false,
		Help:     "Shows all available commands",
	})

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//help",
		Callback: p.adminHelpCommand,
		Admin:    true,
		Help:     "Shows all available admin commands",
	})

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//shutdown",
		Callback: p.shutdownCommand,
		Admin:    true,
		Help:     "Shutdown the controller",
	})

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//reboot",
		Callback: p.rebootCommand,
		Admin:    true,
		Help:     "Reboot the controller",
	})

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//load",
		Callback: p.loadPluginCommand,
		Admin:    true,
		Help:     "Load a plugin",
	})

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//unload",
		Callback: p.unloadPluginCommand,
		Admin:    true,
		Help:     "Unload a plugin",
	})

	return nil
}

func (p *ControllerPlugin) Unload() error {
	return errors.New("Cannot unload controller plugin")
}

func (p *ControllerPlugin) helpCommand(login string, args []string) {
	items := make([][]any, 0)

	for _, command := range app.GetCommandManager().Commands {
		if command.Admin {
			continue
		}

		items = append(items, []any{
			command.Name,
			strings.Join(command.Aliases, " - "),
			command.Help,
		})
	}

	columns := []ui.Column{
		{Name: "Command", Width: 20},
		{Name: "Aliases", Width: 30},
		{Name: "Description", Width: 50},
	}

	window := ui.NewListWindow(&login)
	window.Title = "Available commands"
	window.Columns = columns
	window.Items = items

	go window.Display()
}

func (p *ControllerPlugin) adminHelpCommand(login string, args []string) {
	items := make([][]any, 0)
	for _, command := range app.GetCommandManager().Commands {
		if !command.Admin {
			continue
		}

		items = append(items, []any{
			command.Name,
			strings.Join(command.Aliases, " - "),
			command.Help,
		})
	}

	columns := []ui.Column{
		{Name: "Command", Width: 20},
		{Name: "Aliases", Width: 30},
		{Name: "Description", Width: 50},
	}

	window := ui.NewListWindow(&login)
	window.Title = "Available admin commands"
	window.Columns = columns
	window.Items = items

	go window.Display()
}

func (p *ControllerPlugin) shutdownCommand(login string, args []string) {
	app.GetGoController().Shutdown()
}

func (p *ControllerPlugin) rebootCommand(login string, args []string) {
	if err := app.GetGoController().Reboot(); err != nil {
		go app.GetGoController().Chat("Error rebooting: "+err.Error(), login)
	}
}

func (p *ControllerPlugin) loadPluginCommand(login string, args []string) {
	if len(args) < 1 {
		go app.GetGoController().Chat("Usage: //load [*plugin]", login)
		return
	}

	pluginName := args[0]
	if err := app.GetPluginManager().LoadPlugin(args[0]); err != nil {
		go app.GetGoController().Chat(fmt.Sprintf("Plugin %s not found", pluginName), login)
	} else {
		go app.GetGoController().Chat(fmt.Sprintf("Plugin %s loaded", pluginName), login)
	}
}

func (p *ControllerPlugin) unloadPluginCommand(login string, args []string) {
	if len(args) < 1 {
		go app.GetGoController().Chat("Usage: //unload [*plugin]", login)
		return
	}

	pluginName := args[0]
	if err := app.GetPluginManager().UnloadPlugin(args[0]); err != nil {
		go app.GetGoController().Chat(fmt.Sprintf("Couldn't unload %s: %s", pluginName, err.Error()), login)
	} else {
		go app.GetGoController().Chat(fmt.Sprintf("Plugin %s unloaded", pluginName), login)
	}
}

func init() {
	controllerPlugin := CreateControllerPlugin()
	app.GetPluginManager().PreLoadPlugin(controllerPlugin)
}