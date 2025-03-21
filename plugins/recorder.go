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
	"github.com/MRegterschot/GoController/plugins/windows"
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
	Recording   *database.Recording
}

func CreateRecorderPlugin() *RecorderPlugin {
	return &RecorderPlugin{
		Name:         "Recorder",
		Dependencies: []string{},
		Loaded:       false,
		BasePlugin:   app.GetBasePlugin(),
	}
}

func (p *RecorderPlugin) Load() error {
	commandManager := app.GetCommandManager()

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//recorder",
		Callback: p.recorderCommand,
		Admin:    true,
		Help:     "Start or stop recording",
		Aliases:  []string{"//rec"},
	})

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//export",
		Callback: p.exportToCSVCommand,
		Admin:    true,
		Help:     "Export recording to CSV",
		Aliases:  []string{"//exp"},
	})

	commandManager.AddCommand(models.ChatCommand{
		Name:     "//recordings",
		Callback: p.recordingsCommand,
		Admin:    true,
		Help:     "Show recordings",
		Aliases:  []string{"//recs"},
	})

	return nil
}

func (p *RecorderPlugin) Unload() error {
	commandManager := app.GetCommandManager()

	commandManager.RemoveCommand("//recorder")
	commandManager.RemoveCommand("//export")
	commandManager.RemoveCommand("//recordings")

	return nil
}

func (p *RecorderPlugin) startRecording(name string) {
	mode := p.GoController.MapManager.CurrentMode
	if mode == "Trackmania/TM_TimeAttack_Online.Script.txt" {
		p.GoController.Server.Client.OnPlayerFinish = append(p.GoController.Server.Client.OnPlayerFinish, gbxclient.GbxCallbackStruct[events.PlayerWayPointEventArgs]{
			Key:  "recording",
			Call: p.onPlayerFinish})
		p.createRecording(name, "TimeAttack")
	} else if mode == "Trackmania/TM_Rounds_Online.Script.txt" {
		p.GoController.Server.Client.OnPreEndRound = append(p.GoController.Server.Client.OnPreEndRound, gbxclient.GbxCallbackStruct[events.ScoresEventArgs]{
			Key:  "recording",
			Call: p.onPreEndRound})
		p.createRecording(name, "Rounds")
	} else {
		p.GoController.Server.Client.OnPreEndRound = append(p.GoController.Server.Client.OnPreEndRound, gbxclient.GbxCallbackStruct[events.ScoresEventArgs]{
			Key:  "recording",
			Call: p.onPreEndRoundMatch})
		p.createRecording(name, "Match")
	}

	p.IsRecording = true
}

func (p *RecorderPlugin) createRecording(name string, modeType string) {
	newRecording := database.NewRecording(database.Recording{
		Name: name,
		Mode: p.GoController.MapManager.CurrentMode,
		Type: modeType,
		Maps: []database.MapRecords{},
	})

	_, err := database.InsertRecording(context.Background(), newRecording)
	if err != nil {
		zap.L().Error("Failed to insert recording", zap.Error(err))
		return
	}

	p.Recording = &newRecording
	zap.L().Info("Recording started")
}

func (p *RecorderPlugin) stopRecording() {
	for i, callback := range p.GoController.Server.Client.OnPlayerFinish {
		if callback.Key == "recording" {
			p.GoController.Server.Client.OnPlayerFinish = append(p.GoController.Server.Client.OnPlayerFinish[:i], p.GoController.Server.Client.OnPlayerFinish[i+1:]...)
		}
	}

	for i, callback := range p.GoController.Server.Client.OnPreEndRound {
		if callback.Key == "recording" {
			p.GoController.Server.Client.OnPreEndRound = append(p.GoController.Server.Client.OnPreEndRound[:i], p.GoController.Server.Client.OnPreEndRound[i+1:]...)
		}
	}

	p.IsRecording = false
	zap.L().Info("Recording stopped")
}

func (p *RecorderPlugin) onPlayerFinish(playerFinishEvent events.PlayerWayPointEventArgs) {
	mapId := p.GoController.MapManager.CurrentMapDB.ID

	if len(p.Recording.Maps) == 0 {
		mapRecords := database.MapRecords{
			ID:       primitive.NewObjectID(),
			MapID:    mapId,
			Finishes: []database.PlayerFinish{},
		}

		p.Recording.Maps = append(p.Recording.Maps, mapRecords)
	}

	last := len(p.Recording.Maps) - 1
	if p.Recording.Maps[last].MapID != mapId {
		mapRecords := database.MapRecords{
			ID:       primitive.NewObjectID(),
			MapID:    mapId,
			Finishes: []database.PlayerFinish{},
		}

		p.Recording.Maps = append(p.Recording.Maps, mapRecords)
		last++
	}

	var playerID *primitive.ObjectID
	if playerDB, err := database.GetPlayerByLogin(context.Background(), playerFinishEvent.Login); err != nil {
		zap.L().Error("Player not found", zap.String("login", playerFinishEvent.Login))
	} else {
		playerID = &playerDB.ID
	}

	playerFinish := database.PlayerFinish{
		ID:          primitive.NewObjectID(),
		PlayerID:    playerID,
		Login:       playerFinishEvent.Login,
		AccountId:   playerFinishEvent.AccountId,
		Time:        playerFinishEvent.RaceTime,
		Checkpoints: playerFinishEvent.CurLapCheckpoints,
		Timestamp:   primitive.NewDateTimeFromTime(time.Now()),
	}

	p.Recording.Maps[last].Finishes = append(p.Recording.Maps[last].Finishes, playerFinish)
	p.Recording.Update(*p.Recording)

	_, err := database.UpdateRecording(context.Background(), *p.Recording)
	if err != nil {
		zap.L().Error("Failed to update match recording", zap.Error(err))
		return
	}

	zap.L().Info("Finish recorded")
}

func (p *RecorderPlugin) onPreEndRound(preEndRoundEvent events.ScoresEventArgs) {
	mapId := p.GoController.MapManager.CurrentMapDB.ID

	if len(p.Recording.Maps) == 0 {
		mapRecords := database.MapRecords{
			ID:     primitive.NewObjectID(),
			MapID:  mapId,
			Rounds: []database.Round{},
		}

		p.Recording.Maps = append(p.Recording.Maps, mapRecords)
	}

	last := len(p.Recording.Maps) - 1
	if p.Recording.Maps[last].MapID != mapId {
		mapRecords := database.MapRecords{
			ID:     primitive.NewObjectID(),
			MapID:  mapId,
			Rounds: []database.Round{},
		}

		p.Recording.Maps = append(p.Recording.Maps, mapRecords)
		last++
	}

	var playerRounds []database.PlayerRound
	for _, player := range preEndRoundEvent.Players {
		var playerID *primitive.ObjectID
		if playerDB, err := database.GetPlayerByLogin(context.Background(), player.Login); err != nil {
			zap.L().Error("Player not found", zap.String("login", player.Login))
		} else {
			playerID = &playerDB.ID
		}

		playerRounds = append(playerRounds, database.PlayerRound{
			ID:          primitive.NewObjectID(),
			PlayerID:    playerID,
			Login:       player.Login,
			AccountId:   player.AccountId,
			Points:      player.RoundPoints,
			TotalPoints: player.MapPoints + player.RoundPoints,
			Time:        player.PrevRaceTime,
			Checkpoints: player.PrevRaceCheckpoints,
		})
	}

	round := database.Round{
		ID:          primitive.NewObjectID(),
		RoundNumber: len(p.Recording.Maps[last].Rounds) + 1,
		Players:     playerRounds,
	}

	p.Recording.Maps[last].Rounds = append(p.Recording.Maps[last].Rounds, round)
	p.Recording.Update(*p.Recording)

	_, err := database.UpdateRecording(context.Background(), *p.Recording)
	if err != nil {
		zap.L().Error("Failed to update match recording", zap.Error(err))
		return
	}

	zap.L().Info("Round recorded")
}

func (p *RecorderPlugin) onPreEndRoundMatch(preEndRoundEvent events.ScoresEventArgs) {
	mapId := p.GoController.MapManager.CurrentMapDB.ID

	if len(p.Recording.Maps) == 0 {
		mapRecords := database.MapRecords{
			ID:     primitive.NewObjectID(),
			MapID:  mapId,
			Rounds: []database.Round{},
		}

		p.Recording.Maps = append(p.Recording.Maps, mapRecords)
	}

	last := len(p.Recording.Maps) - 1
	if p.Recording.Maps[last].MapID != mapId {
		mapRecords := database.MapRecords{
			ID:          primitive.NewObjectID(),
			MapID:       mapId,
			MatchRounds: []database.MatchRound{},
		}

		p.Recording.Maps = append(p.Recording.Maps, mapRecords)
		last++
	}

	var teams []database.Team
	for _, team := range preEndRoundEvent.Teams {
		teams = append(teams, database.Team{
			ID:          primitive.NewObjectID(),
			TeamID:      team.ID,
			Name:        team.Name,
			Points:      team.RoundPoints,
			TotalPoints: team.MapPoints,
			Players:     []database.PlayerRound{},
		})
	}

	for _, player := range preEndRoundEvent.Players {
		for i, team := range teams {
			if team.TeamID == player.Team {
				var playerID *primitive.ObjectID
				if playerDB, err := database.GetPlayerByLogin(context.Background(), player.Login); err != nil {
					zap.L().Error("Player not found", zap.String("login", player.Login))
				} else {
					playerID = &playerDB.ID
				}

				teams[i].Players = append(teams[i].Players, database.PlayerRound{
					ID:          primitive.NewObjectID(),
					PlayerID:    playerID,
					Login:       player.Login,
					AccountId:   player.AccountId,
					Points:      player.RoundPoints,
					TotalPoints: player.MapPoints + player.RoundPoints,
					Time:        player.PrevRaceTime,
					Checkpoints: player.PrevRaceCheckpoints,
				})
				break
			}
		}
	}

	matchRound := database.MatchRound{
		ID:          primitive.NewObjectID(),
		RoundNumber: len(p.Recording.Maps[last].MatchRounds) + 1,
		Teams:       teams,
	}

	p.Recording.Maps[last].MatchRounds = append(p.Recording.Maps[last].MatchRounds, matchRound)
	p.Recording.Update(*p.Recording)

	_, err := database.UpdateRecording(context.Background(), *p.Recording)
	if err != nil {
		zap.L().Error("Failed to update match recording", zap.Error(err))
		return
	}

	zap.L().Info("Match round recorded")
}

func (p *RecorderPlugin) recorderCommand(login string, args []string) {
	if len(args) < 1 {
		go p.GoController.Chat("Usage: //recorder [*start | stop] [name]", login)
		return
	}

	switch args[0] {
	case "start":
		if p.IsRecording {
			go p.GoController.Chat("Already recording", login)
			return
		}

		name := time.Now().Format("2006-01-02 15:04:05")
		if len(args) > 1 {
			name = strings.Join(args[1:], " ")
		}

		p.startRecording(name)
		go p.GoController.Chat("Recording started", login)
	case "stop":
		if !p.IsRecording {
			go p.GoController.Chat("Not recording", login)
			return
		}

		p.stopRecording()
		go p.GoController.Chat("Recording stopped with id "+p.Recording.ID.Hex(), login)
	default:
		go p.GoController.Chat("Usage: //recorder [*start | stop] [name]", login)
	}
}

func (p *RecorderPlugin) recordingsCommand(login string, args []string) {
	window := windows.CreateRecorderGridWindow(&login)
	window.Title = "Recordings"
	window.Items = make([]any, 0)

	recordingsDB, err := database.GetRecordings(context.Background())
	if err != nil {
		zap.L().Error("Failed to get recordings", zap.Error(err))
		go p.GoController.Chat("Failed to get recordings", login)
		return
	}

	for _, recordingDB := range recordingsDB {
		var recording models.Recording
		recordingDB.ToModel(&recording)
		window.Actions[recording.ID] = app.GetUIManager().AddAction(p.handleDownloadAnswer, recording.ID)
		window.Items = append(window.Items, recording)
	}
	
	go window.Display()
}

func (p *RecorderPlugin) exportToCSVCommand(login string, args []string) {
	if len(args) < 1 {
		go p.GoController.Chat("Usage: //export [*recording id]", login)
		return
	}

	recordingID := args[0]
	objectID, err := primitive.ObjectIDFromHex(recordingID)
	if err != nil {
		zap.L().Error("Invalid recording ID", zap.Error(err))
		go p.GoController.Chat("Invalid recording ID", login)
		return
	}

	err = p.exportToCSV(objectID)
	if err != nil {
		go p.GoController.Chat("Failed to export recording to CSV", login)
		return
	}

	go p.GoController.Chat("Recording exported to CSV", login)
}

func (p *RecorderPlugin) handleDownloadAnswer(login string, data any, _ any) {
	if id, err := primitive.ObjectIDFromHex(data.(string)); err != nil {
		zap.L().Error("Invalid recording ID", zap.Error(err))
		go p.GoController.Chat("Invalid recording ID", login)
	} else {
		if err = p.exportToCSV(id); err != nil {
			go p.GoController.Chat("Failed to export recording to CSV", login)
		} else {
			go p.GoController.Chat("Recording exported to CSV", login)
		}
	}
}

func (p *RecorderPlugin) exportToCSV(id primitive.ObjectID) error {
	recording, err := database.GetRecordingByID(context.Background(), id)
	if err != nil {
		zap.L().Error("Failed to get recording", zap.Error(err))
		return err
	}

	data := [][]string{
		{"Time", "Track", "PlayerID", "PlayerName", "Record", "RoundNumber", "Checkpoints"},
	}

	switch recording.Type {
	case "TimeAttack":
		for _, mapRecords := range recording.Maps {
			for _, finish := range mapRecords.Finishes {
				checkpoints := strings.Trim(fmt.Sprint(finish.Checkpoints), "[]")
				mapName := "Unknown"
				mapDB, err := database.GetMapByID(context.Background(), mapRecords.MapID)
				if err == nil {
					mapName = mapDB.Name
				}
				data = append(data, []string{
					fmt.Sprint(finish.Timestamp.Time().Unix()),
					mapName,
					finish.AccountId,
					finish.Login,
					fmt.Sprint(finish.Time),
					"",
					checkpoints,
				})
			}
		}
	case "Rounds":
		for _, mapRecords := range recording.Maps {
			for _, round := range mapRecords.Rounds {
				for _, player := range round.Players {
					checkpoints := strings.Trim(fmt.Sprint(player.Checkpoints), "[]")
					mapName := "Unknown"
					mapDB, err := database.GetMapByID(context.Background(), mapRecords.MapID)
					if err == nil {
						mapName = mapDB.Name
					}
					data = append(data, []string{
						"",
						mapName,
						player.AccountId,
						player.Login,
						fmt.Sprint(player.Time),
						fmt.Sprint(round.RoundNumber),
						checkpoints,
					})
				}
			}
		}
	case "Match":
		for _, mapRecords := range recording.Maps {
			for _, round := range mapRecords.MatchRounds {
				for _, team := range round.Teams {
					for _, player := range team.Players {
						checkpoints := strings.Trim(fmt.Sprint(player.Checkpoints), "[]")
						mapName := "Unknown"
						mapDB, err := database.GetMapByID(context.Background(), mapRecords.MapID)
						if err == nil {
							mapName = mapDB.Name
						}
						data = append(data, []string{
							"",
							mapName,
							player.AccountId,
							player.Login,
							fmt.Sprint(player.Time),
							fmt.Sprint(round.RoundNumber),
							checkpoints,
						})
					}
				}
			}
		}
	}

	filePath := "recording_" + id.Hex() + ".csv"
	if err := utils.ExportCSV("./exports/"+filePath, data); err != nil {
		zap.L().Error("Failed to export to CSV", zap.Error(err))
		return err
	}

	return nil
}

func init() {
	recorderPlugin := CreateRecorderPlugin()
	app.GetPluginManager().PreLoadPlugin(recorderPlugin)
}
