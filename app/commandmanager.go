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

		outCommands = append(outCommands, fmt.Sprintf("$0C6%s - %s$FFF %s", command.Name, strings.Join(command.Aliases, " - "), command.Help))
	}

	go GetGoController().Chat("Available commands: "+strings.Join(outCommands, ", "), login)
}

func (cm *CommandManager) adminHelpCommand(login string, args []string) {
	var outCommands []string
	for _, command := range cm.Commands {
		if !command.Admin {
			continue
		}

		outCommands = append(outCommands, fmt.Sprintf("$0C6%s - %s$FFF %s", command.Name, strings.Join(command.Aliases, " - "), command.Help))
	}

	go GetGoController().Chat("Available admin commands: "+strings.Join(outCommands, ", "), login)
}

func (cm *CommandManager) shutdownCommand(login string, args []string) {
	GetGoController().Shutdown()
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
func (cm *CommandManager) ExecuteCommand(login string, text string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if strings.HasPrefix(text, "/") {
		controller := GetGoController()

		for _, command := range cm.Commands {
			if command.Name == "" && len(command.Aliases) == 0 {
				continue
			}

			if strings.HasPrefix(text, "//") && !utils.Includes(*controller.Admins, login) {
				go controller.Chat("$C00Not allowed.", login)
				return
			}

			// Prepare regex
			prefix := `[/]`
			if strings.HasPrefix(command.Name, "//") {
				prefix = `[/]{2}`
			}
			exp := regexp.MustCompile(fmt.Sprintf(`^%s\b%s\b`, prefix, EscapeRegex(strings.TrimLeft(command.Name, "/"))))

			// Match command
			if exp.MatchString(text) {
				// Extract parameters
				words := strings.TrimSpace(strings.Replace(text, command.Name, "", 1))
				params := regexp.MustCompile(`(?:[^"\s]+|"[^"]*")`).FindAllString(words, -1)

				// Remove surrounding quotes
				for i, word := range params {
					params[i] = strings.Trim(word, `"`)
				}

				// Execute command
				go command.Callback(login, params) // Run in a goroutine to mimic async behavior
				zap.L().Debug("Command executed", zap.String("command", command.Name), zap.String("login", login), zap.Strings("params", params))
				return
			} else {
				for _, alias := range command.Aliases {
					exp = regexp.MustCompile(fmt.Sprintf(`^%s\b%s\b`, prefix, EscapeRegex(strings.TrimLeft(alias, "/"))))

					if exp.MatchString(text) {
						// Extract parameters
						words := strings.TrimSpace(strings.Replace(text, alias, "", 1))
						params := regexp.MustCompile(`(?:[^"\s]+|"[^"]*")`).FindAllString(words, -1)

						// Remove surrounding quotes
						for i, word := range params {
							params[i] = strings.Trim(word, `"`)
						}

						// Execute command
						go command.Callback(login, params) // Run in a goroutine to mimic async behavior
						zap.L().Debug("Command executed", zap.String("command", alias), zap.String("login", login), zap.Strings("params", params))
						return
					}
				}
			}
		}
		go controller.Chat(fmt.Sprintf("$fffCommand $0C6%s $fffnot found.", text), login)
	}

}

func (cm *CommandManager) onPlayerChat(chatEvent events.PlayerChatEventArgs) {
	if chatEvent.PlayerUid == 0 {
		return
	}

	cm.ExecuteCommand(chatEvent.Login, chatEvent.Text)
}

// EscapeRegex escapes special regex characters
func EscapeRegex(text string) string {
	replacer := strings.NewReplacer(
		".", `\.`, "*", `\*`, "+", `\+`, "?", `\?`,
		"{", `\{`, "}", `\}`, "(", `\(`, ")", `\)`,
		"[", `\[`, "]", `\]`, "|", `\|`, "^", `\^`,
		"$", `\$`,
	)
	return replacer.Replace(text)
}
