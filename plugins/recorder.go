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
	"github.com/MRegterschot/GoController/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type RecorderPlugin struct {
	app.BasePlugin
	Name         string
	Dependencies []string
	Loaded       bool

	IsRecording bool
	Recording *database.Recording
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
		Callback: m.recorderCommand,
		Admin:    true,
		Help:     "Start or stop recording",
	})

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//export",
		Callback: m.exportToCSVCommand,
		Admin:    true,
		Help:     "Export recording to CSV",
	})

	return nil
}

func (m *RecorderPlugin) Unload() error {
	return nil
}

func (m *RecorderPlugin) startRecording(name string) {
	mode := m.GoController.MapManager.CurrentMode
	if mode == "Trackmania/TM_Rounds_Online.Script.txt" || mode == "Trackmania/TM_TimeAttack_Online.Script.txt" {
		m.GoController.Server.Client.OnPlayerFinish = append(m.GoController.Server.Client.OnPlayerFinish, gbxclient.GbxCallbackStruct[events.PlayerWayPointEventArgs]{
			Key:  "recording",
			Call: m.onPlayerFinish})
		m.createRecording(name, "Rounds")
		} else {
		m.GoController.Server.Client.OnEndRound = append(m.GoController.Server.Client.OnEndRound, gbxclient.GbxCallbackStruct[events.ScoresEventArgs]{
			Key:  "recording",
			Call: m.onEndRound})
		m.createRecording(name, "Match")
	}

	m.IsRecording = true
}

func (m *RecorderPlugin) createRecording(name string, modeType string) {
	newRecording := database.NewRecording(database.Recording{
		Name: name,
		Mode: m.GoController.MapManager.CurrentMode,
		Type: modeType,
		Maps: []database.MapRecords{},
	})

	_, err := database.InsertRecording(context.Background(), newRecording)
	if err != nil {
		zap.L().Error("Failed to insert recording", zap.Error(err))
		return
	}

	m.Recording = &newRecording
	zap.L().Info("Recording started")
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
	mapId := m.GoController.MapManager.CurrentMapDB.ID

	if len(m.Recording.Maps) == 0 {
		mapRecords := database.MapRecords{
			ID: primitive.NewObjectID(),
			MapID: mapId,
			Rounds: []database.PlayerRound{},
		}
	
		m.Recording.Maps = append(m.Recording.Maps, mapRecords)
	}

	last := len(m.Recording.Maps) - 1
	if m.Recording.Maps[last].MapID != mapId {
		mapRecords := database.MapRecords{
			ID: primitive.NewObjectID(),
			MapID: mapId,
			Rounds: []database.PlayerRound{},
		}

		m.Recording.Maps = append(m.Recording.Maps, mapRecords)
		last++
	}

	var playerID *primitive.ObjectID
	if playerDB, err := database.GetPlayerByLogin(context.Background(), playerFinishEvent.Login); err != nil {
		zap.L().Error("Player not found", zap.String("login", playerFinishEvent.Login))
	} else {
		playerID = &playerDB.ID
	}

	round := database.PlayerRound{
		ID: primitive.NewObjectID(),
		PlayerID: playerID,
		Login: playerFinishEvent.Login,
		AccountId: playerFinishEvent.AccountId,
		Time: playerFinishEvent.RaceTime,
		Checkpoints: playerFinishEvent.CurLapCheckpoints,
		Timestamp: primitive.NewDateTimeFromTime(time.Now()),
	}

	m.Recording.Maps[last].Rounds = append(m.Recording.Maps[last].Rounds, round)
	m.Recording.Update(*m.Recording)

	_, err := database.UpdateRecording(context.Background(), *m.Recording)
	if err != nil {
		zap.L().Error("Failed to update match recording", zap.Error(err))
		return
	}

	zap.L().Info("Round recorded")
}

func (m *RecorderPlugin) onEndRound(_ *gbxclient.GbxClient, endRoundEvent events.ScoresEventArgs) {
	mapId := m.GoController.MapManager.CurrentMapDB.ID

	if len(m.Recording.Maps) == 0 {
		mapRecords := database.MapRecords{
			ID: primitive.NewObjectID(),
			MapID: mapId,
			MatchRounds: []database.MatchRound{},
		}
	
		m.Recording.Maps = append(m.Recording.Maps, mapRecords)
	}

	last := len(m.Recording.Maps) - 1
	if m.Recording.Maps[last].MapID != mapId {
		mapRecords := database.MapRecords{
			ID: primitive.NewObjectID(),
			MapID: mapId,
			MatchRounds: []database.MatchRound{},
		}

		m.Recording.Maps = append(m.Recording.Maps, mapRecords)
		last++
	}
	
	var teams []database.Team
	for _, team := range endRoundEvent.Teams {
		teams = append(teams, database.Team{
			ID: primitive.NewObjectID(),
			TeamID: team.ID,
			Name: team.Name,
			Points: team.RoundPoints,
			TotalPoints: team.MapPoints,
			Players: []database.PlayerMatchRound{},
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

				teams[i].Players = append(teams[i].Players, database.PlayerMatchRound{
					ID: primitive.NewObjectID(),
					PlayerID: playerID,
					Login: player.Login,
					AccountId: player.AccountId,
					Points: player.RoundPoints,
					TotalPoints: player.MapPoints,
					Time: player.PrevRaceTime,
					Checkpoints: player.PrevRaceCheckpoints,
				})
				break
			}
		}
	}

	round := database.MatchRound{
		ID: primitive.NewObjectID(),
		RoundNumber: len(m.Recording.Maps[last].Rounds) + 1,
		Teams: teams,
	}

	m.Recording.Maps[last].MatchRounds = append(m.Recording.Maps[last].MatchRounds, round)
	m.Recording.Update(*m.Recording)

	_, err := database.UpdateRecording(context.Background(), *m.Recording)
	if err != nil {
		zap.L().Error("Failed to update match recording", zap.Error(err))
		return
	}

	zap.L().Info("Round recorded")
}

func (m *RecorderPlugin) recorderCommand(login string, args []string) {
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
		go m.GoController.Chat("Recording stopped with id " + m.Recording.ID.Hex(), login)
	default:
		go m.GoController.Chat("Usage: //recorder [start | stop]", login)
	}
}

func (m *RecorderPlugin) exportToCSVCommand(login string, args []string) {
	if len(args) < 1 {
		go m.GoController.Chat("Usage: //export [*recording id]", login)
		return
	}

	recordingID := args[0]
	objectID, err := primitive.ObjectIDFromHex(recordingID)
	if err != nil {
		zap.L().Error("Invalid recording ID", zap.Error(err))
		go m.GoController.Chat("Invalid recording ID", login)
		return
	}
	recording, err := database.GetRecordingByID(context.Background(), objectID)
	if err != nil {
		zap.L().Error("Failed to get recording", zap.Error(err))
		go m.GoController.Chat("Failed to get recording", login)
		return
	}

	data := [][]string{
		{"Time", "Track", "PlayerID", "PlayerName", "Record", "RoundNumber", "Checkpoints"},
	}
	for _, mapRecords := range recording.Maps {
		for _, round := range mapRecords.Rounds {
			checkpoints := strings.Trim(fmt.Sprint(round.Checkpoints), "[]")
			data = append(data, []string{
				fmt.Sprint(round.Timestamp.Time().Unix()),
				mapRecords.MapID.Hex(),
				round.AccountId,
				round.Login,
				fmt.Sprint(round.Time),
				"",
				checkpoints,
			})
		}
	}

	filePath := "recording_" + recordingID + ".csv"
	if err := utils.ExportCSV("./exports/" + filePath, data); err != nil {
		zap.L().Error("Failed to export to CSV", zap.Error(err))
		go m.GoController.Chat("Failed to export to CSV", login)
		return
	}

	go m.GoController.Chat("Exported to " + filePath, login)
}

func init() {
	recorderPlugin := CreateRecorderPlugin()
	app.GetPluginManager().PreLoadPlugin(recorderPlugin)
}
