package app

import (
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

	GetClient().OnPlayerChat = append(GetClient().OnPlayerChat, gbxclient.GbxCallbackStruct[events.PlayerChatEventArgs]{
		Key:  "cmPlayerChat",
		Call: cm.onPlayerChat})

	zap.L().Info("CommandManager initialized")
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

	if admin && !controller.IsAdmin(login) {
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
	if !utils.CommandRegex.MatchString(chatEvent.Text) {
		return
	}

	admin := false
	if strings.HasPrefix(chatEvent.Text, "//") {
		admin = true
	}
		
	splitText := strings.Split(chatEvent.Text, " ")
	cm.ExecuteCommand(chatEvent.Login, splitText[0], splitText[1:], admin)
}
