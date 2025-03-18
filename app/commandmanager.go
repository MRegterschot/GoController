package app

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/MRegterschot/GbxRemoteGo/events"
	"github.com/MRegterschot/GbxRemoteGo/gbxclient"
	. "github.com/MRegterschot/GoController/models"
	"github.com/MRegterschot/GoController/utils"
	"go.uber.org/zap"
)

type CommandManager struct {
	Commands map[string]ChatCommand
	mu       sync.Mutex
}

var (
	cmInstance *CommandManager
	cmOnce     sync.Once
)

var re = regexp.MustCompile(`^/{1,2}`)

func GetCommandManager() *CommandManager {
	cmOnce.Do(func() {
		cmInstance = &CommandManager{
			Commands: make(map[string]ChatCommand),
		}
	})
	return cmInstance
}

func (cm *CommandManager) Init() {
	zap.L().Info("Initializing CommandManager")

	cm.addDefaultCommands()

	GetClient().OnPlayerChat = append(GetClient().OnPlayerChat, gbxclient.GbxCallbackStruct[events.PlayerChatEventArgs]{
		Key:  "cmPlayerChat",
		Call: cm.onPlayerChat})

	zap.L().Info("CommandManager initialized")
}

// Adds the default commands to the CommandManager
func (cm *CommandManager) addDefaultCommands() {
	cm.AddCommand(ChatCommand{
		Name:     "/help",
		Callback: cm.helpCommand,
		Admin:    false,
		Help:     "Shows all available commands",
	})

	cm.AddCommand(ChatCommand{
		Name:     "//help",
		Callback: cm.adminHelpCommand,
		Admin:    true,
		Help:     "Shows all available admin commands",
	})

	cm.AddCommand(ChatCommand{
		Name:     "//shutdown",
		Callback: cm.shutdownCommand,
		Admin:    true,
		Help:     "Shuts down the controller",
	})

	cm.AddCommand(ChatCommand{
		Name:     "//reboot",
		Callback: cm.rebootCommand,
		Admin:    true,
		Help:     "Reboots the controller",
	})

	cm.AddCommand(ChatCommand{
		Name:     "//load",
		Callback: cm.loadPluginCommand,
		Admin:    true,
		Help:     "Loads a plugin",
	})

	cm.AddCommand(ChatCommand{
		Name:     "//unload",
		Callback: cm.unloadPluginCommand,
		Admin:    true,
		Help:     "Unloads a plugin",
	})
}

// The default commands
func (cm *CommandManager) helpCommand(login string, args []string) {
	var outCommands []string

	for _, command := range cm.Commands {
		if command.Admin {
			continue
		}

		msg := "$0C6" + command.Name
		if len(command.Aliases) > 0 {
			msg += " - " + strings.Join(command.Aliases, " - ")
		}
		msg += "$FFF " + command.Help
		outCommands = append(outCommands, msg)
	}

	go GetGoController().Chat("Available commands: "+strings.Join(outCommands, ", "), login)
}

func (cm *CommandManager) adminHelpCommand(login string, args []string) {
	var outCommands []string
	for _, command := range cm.Commands {
		if !command.Admin {
			continue
		}
		msg := "$0C6" + command.Name
		if len(command.Aliases) > 0 {
			msg += " - " + strings.Join(command.Aliases, " - ")
		}
		msg += "$FFF " + command.Help
		outCommands = append(outCommands, msg)
	}

	go GetGoController().Chat("Available admin commands: "+strings.Join(outCommands, ", "), login)
}

func (cm *CommandManager) shutdownCommand(login string, args []string) {
	GetGoController().Shutdown()
}

func (cm *CommandManager) rebootCommand(login string, args []string) {
	if err := GetGoController().Reboot(); err != nil {
		go GetGoController().Chat("Failed to reboot", login)
	}
}

func (cm *CommandManager) loadPluginCommand(login string, args []string) {
	if len(args) < 1 {
		go GetGoController().Chat("Usage: //loadplugin [*plugin]", login)
		return
	}

	pluginName := args[0]
	if err := GetPluginManager().LoadPlugin(args[0]); err != nil {
		go GetGoController().Chat(fmt.Sprintf("Plugin %s not found", pluginName), login)
	} else {
		go GetGoController().Chat(fmt.Sprintf("Plugin %s loaded", pluginName), login)
	}
}

func (cm *CommandManager) unloadPluginCommand(login string, args []string) {
	if len(args) < 1 {
		go GetGoController().Chat("Usage: //unloadplugin [*plugin]", login)
		return
	}

	pluginName := args[0]
	if err := GetPluginManager().UnloadPlugin(args[0]); err != nil {
		go GetGoController().Chat(fmt.Sprintf("Plugin %s not found", pluginName), login)
	} else {
		go GetGoController().Chat(fmt.Sprintf("Plugin %s unloaded", pluginName), login)
	}
}

// Adds a command to the CommandManager
func (cm *CommandManager) AddCommand(command ChatCommand) {
	cm.Commands[command.Name] = command
	zap.L().Debug("Command added", zap.String("command", command.Name))
}

// Removes a command from the CommandManager
func (cm *CommandManager) RemoveCommand(command string) {
	delete(cm.Commands, command)
	zap.L().Debug("Command removed", zap.String("command", command))
}

// Executes a command
func (cm *CommandManager) ExecuteCommand(login string, command string, params []string, admin bool) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	controller := GetGoController()

	if admin && !utils.Includes(*controller.Admins, login) {
		return
	}

	for _, com := range cm.Commands {
		if com.Name == "" && len(com.Aliases) == 0 {
			continue
		}

		if command == com.Name {
			go com.Callback(login, params)
			zap.L().Debug("Command executed", zap.String("command", com.Name), zap.String("login", login), zap.Strings("params", params))
			return
		} else {
			for _, alias := range com.Aliases {
				if command == alias {
					go com.Callback(login, params)
					zap.L().Debug("Command executed", zap.String("command", alias), zap.String("login", login), zap.Strings("params", params))
					return
				}
			}
		}
	}
}

func (cm *CommandManager) onPlayerChat(chatEvent events.PlayerChatEventArgs) {
	// Check if the player is the server
	if chatEvent.PlayerUid == 0 {
		return
	}

	// Check if the message is a command
	if !re.MatchString(chatEvent.Text) {
		return
	}

	admin := false
	if strings.HasPrefix(chatEvent.Text, "//") {
		admin = true
	}
		
	splitText := strings.Split(chatEvent.Text, " ")
	cm.ExecuteCommand(chatEvent.Login, splitText[0], splitText[1:], admin)
}
