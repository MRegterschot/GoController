package plugins

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/MRegterschot/GbxRemoteGo/events"
	"github.com/MRegterschot/GbxRemoteGo/gbxclient"
	"github.com/MRegterschot/GoController/app"
	"github.com/MRegterschot/GoController/database"
	"github.com/MRegterschot/GoController/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type RecorderPlugin struct {
	app.BasePlugin
	Name         string
	Dependencies []string
	Loaded       bool

	IsRecording bool
	MatchRecording *database.MatchRecording
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

func (m *RecorderPlugin) startRecording(name string) {
	mode := m.GoController.MapManager.CurrentMode
	if mode == "Trackmania/TM_Rounds_Online.Script.txt" || mode == "Trackmania/TM_TimeAttack_Online.Script.txt" {
		m.startRoundsRecording(name)
	} else {
		m.startMatchRecording(name)
	}

	m.IsRecording = true
}

func (m *RecorderPlugin) startMatchRecording(name string) {
	m.GoController.Server.Client.OnEndRound = append(m.GoController.Server.Client.OnEndRound, gbxclient.GbxCallbackStruct[events.ScoresEventArgs]{
		Key:  "recording",
		Call: m.onEndRound})

	newMatchRecording := database.NewMatchRecording(database.MatchRecording{
		Name: name,
		Mode: m.GoController.MapManager.CurrentMode,
		Maps: []database.MapRecords{},
	})

	_, err := database.InsertMatchRecording(context.Background(), newMatchRecording)
	if err != nil {
		zap.L().Error("Failed to insert match recording", zap.Error(err))
		return
	}

	m.MatchRecording = &newMatchRecording
	zap.L().Info("Match recording started")
}

func (m *RecorderPlugin) startRoundsRecording(name string) {
	m.GoController.Server.Client.OnPlayerFinish = append(m.GoController.Server.Client.OnPlayerFinish, gbxclient.GbxCallbackStruct[events.PlayerWayPointEventArgs]{
		Key:  "recording",
		Call: m.onPlayerFinish})
}

func (m *RecorderPlugin) stopRecording() {
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
	mapId := m.GoController.MapManager.CurrentMapDB.ID

	recordingContainsMap := false
	for _, mapRecords := range m.MatchRecording.Maps {
		if mapRecords.MapID == mapId {
			recordingContainsMap = true
			break
		}
	}

	if !recordingContainsMap {
		mapRecords := database.MapRecords{
			ID: primitive.NewObjectID(),
			MapID: mapId,
			Rounds: []database.Round{},
		}

		m.MatchRecording.Maps = append(m.MatchRecording.Maps, mapRecords)
	}

	index := len(m.MatchRecording.Maps) - 1
	roundNumber := len(m.MatchRecording.Maps[index].Rounds) + 1
	
	var teams []database.Team
	for _, team := range endRoundEvent.Teams {
		teams = append(teams, database.Team{
			ID: primitive.NewObjectID(),
			TeamID: team.ID,
			Name: team.Name,
			Points: team.RoundPoints,
			TotalPoints: team.MapPoints,
			Players: []database.PlayerRound{},
		})
	}

	for _, player := range endRoundEvent.Players {
		for i, team := range teams {
			if team.TeamID == player.Team {
				var playerID *primitive.ObjectID
				if playerDB, err := database.GetPlayerByLogin(context.Background(), player.Login); err != nil {
					zap.L().Error("Player not found", zap.String("login", player.Login))
				} else {
					playerID = &playerDB.ID
				}

				teams[i].Players = append(teams[i].Players, database.PlayerRound{
					ID: primitive.NewObjectID(),
					PlayerID: playerID,
					Login: player.Login,
					Points: player.RoundPoints,
					TotalPoints: player.MapPoints,
					Time: player.PrevRaceTime,
					Checkpoints: player.PrevRaceCheckpoints,
				})
				break
			}
		}
	}

	round := database.Round{
		ID: primitive.NewObjectID(),
		RoundNumber: roundNumber,
		Teams: teams,
	}

	m.MatchRecording.Maps[index].Rounds = append(m.MatchRecording.Maps[index].Rounds, round)
	m.MatchRecording.Update(*m.MatchRecording)

	_, err := database.UpdateMatchRecording(context.Background(), *m.MatchRecording)
	if err != nil {
		zap.L().Error("Failed to update match recording", zap.Error(err))
		return
	}

	zap.L().Info("Round recorded")
}

func (m *RecorderPlugin) RecorderCommand(login string, args []string) {
	if len(args) < 1 {
		go m.GoController.Chat("Usage: //recorder [start | stop] [*name]", login)
		return
	}

	
	switch args[0] {
	case "start":
		if m.IsRecording {
			go m.GoController.Chat("Already recording", login)
			return
		}

		name := time.Now().Format("2006-01-02 15:04:05")
		if len(args) > 1 {
			name = strings.Join(args[1:], " ")
		}

		m.startRecording(name)
		go m.GoController.Chat("Recording started", login)
	case "stop":
		if !m.IsRecording {
			go m.GoController.Chat("Not recording", login)
			return
		}

		m.stopRecording()
		go m.GoController.Chat("Recording stopped", login)
	default:
		go m.GoController.Chat("Usage: //recorder [start | stop]", login)
	}
}

func init() {
	recorderPlugin := CreateRecorderPlugin()
	app.GetPluginManager().PreLoadPlugin(recorderPlugin)
}
