package plugins

import (
	"fmt"

	"github.com/MRegterschot/GbxRemoteGo/events"
	"github.com/MRegterschot/GbxRemoteGo/gbxclient"
	"github.com/MRegterschot/GoController/app"
	"github.com/MRegterschot/GoController/models"
	"go.uber.org/zap"
)

type RecorderPlugin struct {
	app.BasePlugin
	Name         string
	Dependencies []string
	Loaded       bool

	IsRecording bool
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
	mode := m.GoController.MapManager.CurrentMode
	if mode == "Trackmania/TM_Rounds_Online.Script.txt" || mode == "Trackmania/TM_TimeAttack_Online.Script.txt" {
		m.GoController.Server.Client.OnPlayerFinish = append(m.GoController.Server.Client.OnPlayerFinish, gbxclient.GbxCallbackStruct[events.PlayerWayPointEventArgs]{
			Key:  "recording",
			Call: m.onPlayerFinish})
	} else {
		m.GoController.Server.Client.OnEndRound = append(m.GoController.Server.Client.OnEndRound, gbxclient.GbxCallbackStruct[events.ScoresEventArgs]{
			Key:  "recording",
			Call: m.onEndRound})
	}

	m.IsRecording = true
	zap.L().Info("Recording started")
}

func (m *RecorderPlugin) StopRecording() {
	for i, callback := range m.GoController.Server.Client.OnPlayerFinish {
		if callback.Key == "recording" {
			m.GoController.Server.Client.OnPlayerFinish = append(m.GoController.Server.Client.OnPlayerFinish[:i], m.GoController.Server.Client.OnPlayerFinish[i+1:]...)
		}
	}

	for i, callback := range m.GoController.Server.Client.OnEndRound {
		if callback.Key == "recording" {
			m.GoController.Server.Client.OnEndRound = append(m.GoController.Server.Client.OnEndRound[:i], m.GoController.Server.Client.OnEndRound[i+1:]...)
		}
	}

	m.IsRecording = false	
	zap.L().Info("Recording stopped")
}

func (m *RecorderPlugin) onPlayerFinish(_ *gbxclient.GbxClient, playerFinishEvent events.PlayerWayPointEventArgs) {
	fmt.Println(playerFinishEvent)
}

func (m *RecorderPlugin) onEndRound(_ *gbxclient.GbxClient, endRoundEvent events.ScoresEventArgs) {
	fmt.Println(endRoundEvent)
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
