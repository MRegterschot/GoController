package app

import (
	"fmt"
	"os"
	"strings"
	"sync"

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
	PluginManager   *PluginManager
	PlayerManager   *PlayerManager
	DatabaseManager *DatabaseManager
	MapManager      *MapManager
}

var (
	gcInstance *GoController
	gcOnce     sync.Once
)

func GetGoController() *GoController {
	gcOnce.Do(func() {
		commandManager := GetCommandManager()
		settingsManager := GetSettingsManager()
		pluginManager := GetPluginManager()
		playerManager := GetPlayerManager()
		databaseManager := GetDatabaseManager()
		mapManager := GetMapManager()

		gcInstance = &GoController{
			StartTime:       utils.GetCurrentTimeInMilliseconds(),
			Version:         "1.0.0",
			Server:          NewServer(),
			CommandManager:  commandManager,
			SettingsManager: settingsManager,
			PluginManager:   pluginManager,
			PlayerManager:   playerManager,
			DatabaseManager: databaseManager,
			MapManager:      mapManager,
			Admins:          &settingsManager.Admins,
		}
	})
	return gcInstance
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

	c.SettingsManager.Init()
	c.DatabaseManager.Init()
	c.CommandManager.Init()
	c.PluginManager.Init()
	c.PlayerManager.Init()
	c.MapManager.Init()

	c.Server.Client.Echo(fmt.Sprintf("%d", c.StartTime), "GoController")

	msg := fmt.Sprintf("Welcome to $0C6GoController$FFF! Version %s", c.Version)
	go c.Chat(msg)
	zap.L().Info(msg)
	zap.L().Info("GoController started successfully")
}

// Sends a chat message to the server
func (c *GoController) Chat(message string, login ...string) {
	if len(login) > 0 {
		go c.Server.Client.ChatSendServerMessageToLogin("$9ab$n>$z$s "+message, strings.Join(login, ","))
	} else {
		go c.Server.Client.ChatSendServerMessage("$9ab»$z$s " + message)
	}
}

// Checks if a login is an admin
func (c *GoController) IsAdmin(login string) bool {
	for _, admin := range *c.Admins {
		if admin == login {
			return true
		}
	}
	return false
}

// Shutdown the GoController
func (c *GoController) Shutdown() {
	zap.L().Info("Shutting down GoController")
	go c.Chat("GoController shutting down...")

	c.Server.Disconnect()

	zap.L().Info("GoController shutdown")
	os.Exit(0)
}
