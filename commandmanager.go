package main

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/MRegterschot/GbxRemoteGo/events"
	"github.com/MRegterschot/GbxRemoteGo/gbxclient"
	"go.uber.org/zap"
)

type CommandCallback func(login string, args []string)

type ChatCommand struct {
	Name     string
	Callback CommandCallback
	Admin    bool
	Help     string
}

type CommandManager struct {
	Commands map[string]ChatCommand
	mu       sync.Mutex
}

// Creates a new CommandManager
func NewCommandManager() *CommandManager {
	return &CommandManager{
		Commands: make(map[string]ChatCommand),
	}
}

func (cm *CommandManager) Init() {
	zap.L().Info("Initializing CommandManager")
	cm.addDefaultCommands()
	GetController().Server.Client.OnPlayerChat = append(GetController().Server.Client.OnPlayerChat, cm.onPlayerChat)
	zap.L().Info("CommandManager initialized")
}

// Adds the default commands to the CommandManager
func (cm *CommandManager) addDefaultCommands() {
	cm.AddCommand(ChatCommand{
		Name:     "/help",
		Callback: cm.HelpCommand,
		Admin:    false,
		Help:     "Shows all available commands",
	})

	cm.AddCommand(ChatCommand{
		Name:     "//help",
		Callback: cm.AdminHelpCommand,
		Admin:    true,
		Help:     "Shows all available admin commands",
	})
}

// The default commands
func (cm *CommandManager) HelpCommand(login string, args []string) {
	var outCommands []string

	for _, command := range cm.Commands {
		if command.Admin {
			continue
		}

		outCommands = append(outCommands, fmt.Sprintf("$0C6%s$FFF %s", command.Name, command.Help))
	}

	GetController().Chat("Available commands: "+strings.Join(outCommands, ", "), login)
}

func (cm *CommandManager) AdminHelpCommand(login string, args []string) {
	var outCommands []string

	for _, command := range cm.Commands {
		if !command.Admin {
			continue
		}

		outCommands = append(outCommands, fmt.Sprintf("$0C6%s$FFF %s", command.Name, command.Help))
	}

	GetController().Chat("Available admin commands: "+strings.Join(outCommands, ", "), login)
}

// Adds a command to the CommandManager
func (cm *CommandManager) AddCommand(command ChatCommand) {
	cm.Commands[command.Name] = command
}

// Removes a command from the CommandManager
func (cm *CommandManager) RemoveCommand(command string) {
	delete(cm.Commands, command)
}

// Executes a command
func (cm *CommandManager) ExecuteCommand(login string, text string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if strings.HasPrefix(text, "/") {
		for _, command := range cm.Commands {
			if command.Name == "" {
				continue
			}

			if strings.HasPrefix(text, "//") {
				return
			}

			// Prepare regex
			prefix := `[/]`
			if strings.HasPrefix(command.Name, "//") {
				prefix = `[/]{2}`
			}
			exp := regexp.MustCompile(fmt.Sprintf(`^%s\b%s\b`, prefix, EscapeRegex(strings.TrimPrefix(command.Name, "/"))))

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
				return
			}
		}
		GetController().Chat(fmt.Sprintf("$fffCommand ¤cmd¤%s $fffnot found.", text), login)
	}

}

func (cm *CommandManager) onPlayerChat(client *gbxclient.GbxClient, chatEvent events.PlayerChatEventArgs) {
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
