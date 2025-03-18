package app

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"

	"github.com/MRegterschot/GbxRemoteGo/gbxclient"
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
	UIManager       *UIManager
}

var (
	gcInstance *GoController
	gcOnce     sync.Once
)

var (
	cInstance *gbxclient.GbxClient
	cOnce     sync.Once
)

func GetGoController() *GoController {
	gcOnce.Do(func() {
		commandManager := GetCommandManager()
		settingsManager := GetSettingsManager()
		pluginManager := GetPluginManager()
		playerManager := GetPlayerManager()
		databaseManager := GetDatabaseManager()
		mapManager := GetMapManager()
		uiManager := GetUIManager()

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
			UIManager:       uiManager,
			Admins:          &settingsManager.Admins,
		}
	})
	return gcInstance
}

func GetClient() *gbxclient.GbxClient {
	cOnce.Do(func() {
		cInstance = GetGoController().Server.Client
	})
	return cInstance
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
	c.Server.Client.TriggerModeScriptEventArray("XmlRpc.EnableCallbacks", []string{"true"})

	c.SettingsManager.Init()
	c.DatabaseManager.Init()
	c.UIManager.Init()
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

func (c *GoController) AfterStart() {
	c.UIManager.AfterInit()
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

// Reboot the GoController
func (c *GoController) Reboot() error {
	exe, err := os.Executable()
	if err != nil {
		zap.L().Error("Failed to get executable path", zap.Error(err))
		return err
	}

	zap.L().Info("Rebooting GoController")
	go c.Chat("GoController rebooting...")

	args := os.Args
	var cmd *exec.Cmd
	if len(args) > 1 {
		cmd = exec.Command(exe, args[1:]...)
	} else {
		cmd = exec.Command(exe)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}

	if err := cmd.Start(); err != nil {
		zap.L().Error("Failed to start new process", zap.Error(err))
		return err
	}

	zap.L().Info("Started new process", zap.String("pid", fmt.Sprintf("%d", cmd.Process.Pid)))

	os.Exit(0)
	return nil
}