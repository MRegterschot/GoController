package main

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/MRegterschot/GoController/plugins"
	"github.com/MRegterschot/GoController/utils"
	"go.uber.org/zap"
)

type GoController struct {
	StartTime       int
	Version         string
	Server          *Server
	MapsPath        string
	Admins          *[]string
	CommandManager  *CommandManager
	SettingsManager *SettingsManager
	PluginManager   *plugins.PluginManager
}

var (
	instance *GoController
	once     sync.Once
)

func GetController() *GoController {
	once.Do(func() {
		commandManager := NewCommandManager()
		settingsManager := NewSettingsManager()
		pluginManager := plugins.GetPluginManager()

		instance = &GoController{
			StartTime:       utils.GetCurrentTimeInMilliseconds(),
			Version:         "1.0.0",
			Server:          NewServer(),
			CommandManager:  commandManager,
			SettingsManager: settingsManager,
			PluginManager:   pluginManager,
			Admins:          &settingsManager.Admins,
		}
	})
	return instance
}

// Starts the GoController
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
	c.PluginManager.Init()

	c.Server.Client.Echo(fmt.Sprintf("%d", c.StartTime), "GoController")

	msg := fmt.Sprintf("Welcome to $0C6GoController$FFF! Version %s", c.Version)
	c.Chat(msg)
	zap.L().Info(msg)
	zap.L().Info("GoController started successfully")
}

// Sends a chat message to the server
func (c *GoController) Chat(message string, login ...string) {
	if len(login) > 0 {
		c.Server.Client.ChatSendServerMessageToLogin("$9ab$n>$z$s "+message, strings.Join(login, ","))
	} else {
		c.Server.Client.ChatSendServerMessage("$9ab»$z$s " + message)
	}
}

// Shutdown the GoController
func (c *GoController) Shutdown() {
	zap.L().Info("Shutting down GoController")
	c.Chat("GoController shutting down...")

	c.Server.Disconnect()

	zap.L().Info("GoController shutdown")
	os.Exit(0)
}
