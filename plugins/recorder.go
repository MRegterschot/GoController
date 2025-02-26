package plugins

import (
	"fmt"

	"github.com/MRegterschot/GbxRemoteGo/events"
	"github.com/MRegterschot/GbxRemoteGo/gbxclient"
	"github.com/MRegterschot/GoController/app"
	"github.com/MRegterschot/GoController/models"
)

type RecorderPlugin struct {
	app.BasePlugin
	Name         string
	Dependencies []string
	Loaded       bool
}

func CreateRecorderPlugin() *RecorderPlugin {
	return &RecorderPlugin{
		Name:         "Recorder",
		Dependencies: []string{},
		Loaded:       false,
		BasePlugin:   app.GetBasePlugin(),
	}
}

func (m *RecorderPlugin) Load() error {
	commandManager := app.GetCommandManager()
	
	m.GoController.Server.Client.OnPlayerFinish = append(m.GoController.Server.Client.OnPlayerFinish, m.onPlayerFinish)
	m.GoController.Server.Client.OnAnyCallback = append(m.GoController.Server.Client.OnAnyCallback, m.onAnyCallback)
	commandManager.AddCommand(models.ChatCommand{
		Name:     "//recorder",
		Callback: m.RecorderCommand,
		Admin:    true,
		Help:     "Start or stop recording",
	})

	return nil
}

func (m *RecorderPlugin) Unload() error {
	return nil
}

func (m *RecorderPlugin) StartRecording() {
	fmt.Println("Recording started")
}

func (m *RecorderPlugin) StopRecording() {
	fmt.Println("Recording stopped")
}

func (m *RecorderPlugin) onPlayerFinish(_ *gbxclient.GbxClient, playerFinishEvent events.PlayerFinishEventArgs) {
	fmt.Println(playerFinishEvent)
}

func (m *RecorderPlugin) onAnyCallback(_ *gbxclient.GbxClient, anyCallbackEvent gbxclient.CallbackEventArgs) {
	fmt.Println(anyCallbackEvent)
}

func (m *RecorderPlugin) RecorderCommand(login string, args []string) {
	if len(args) < 1 {
		go m.GoController.Chat("Usage: //recorder [start | stop]", login)
		return
	}

	switch args[0] {
	case "start":
		m.StartRecording()
		go m.GoController.Chat("Recording started", login)
	case "stop":
		m.StopRecording()
		go m.GoController.Chat("Recording stopped", login)
	default:
		go m.GoController.Chat("Usage: //recorder [start | stop]", login)
	}
}

func init() {
	recorderPlugin := CreateRecorderPlugin()
	app.GetPluginManager().PreLoadPlugin(recorderPlugin)
}
