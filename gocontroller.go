package main

import (
	"strings"
	"sync"

	"go.uber.org/zap"
)

type GoController struct {
	Server          *Server
	MapsPath        string
	Admins          []string
	CommandManager  *CommandManager
	SettingsManager *SettingsManager
}

var (
	instance *GoController
	once     sync.Once
)

func GetController() *GoController {
	once.Do(func() {
		instance = &GoController{
			Server:          NewServer(),
			CommandManager:  NewCommandManager(),
			SettingsManager: NewSettingsManager(),
		}
	})
	return instance
}

func (c *GoController) Start() {
	zap.L().Info("Starting GoController")

	if err := c.Server.Connect(); err != nil {
		zap.L().Fatal("Failed to connect to server", zap.Error(err))
	}

	if err := c.Server.Authenticate(); err != nil {
		zap.L().Fatal("Failed to authenticate with server", zap.Error(err))
	}

	c.Server.Client.EnableCallbacks(true)
	c.Server.Client.SendHideManialinkPage()

	c.Server.Client.SetApiVersion("2023-04-16")
	if mapsPath, err := c.Server.Client.GetMapsDirectory(); err != nil {
		zap.L().Fatal("Failed to get maps directory", zap.Error(err))
	} else {
		c.MapsPath = mapsPath
	}
	c.Server.Client.TriggerModeScriptEvent("XmlRpc.EnableCallbacks", "true")

	c.CommandManager.Init()
	c.SettingsManager.Init()

	zap.L().Info("GoController started")
}

// Sends a chat message to the server
func (c *GoController) Chat(message string, login ...string) {
	if len(login) > 0 {
		c.Server.Client.ChatSendServerMessageToLogin("$9ab$n>$z$s "+message, strings.Join(login, ","))
	} else {
		c.Server.Client.ChatSendServerMessage("$9ab»$z$s " + message)
	}
}
