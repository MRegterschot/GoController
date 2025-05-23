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
	"github.com/MRegterschot/GoController/plugins/widgets"
	"github.com/MRegterschot/GoController/ui"
	"github.com/MRegterschot/GoController/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"slices"
)

type RecorderPlugin struct {
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

	acw := widgets.GetAdminControlsWidget()

	acw.AddAction(widgets.Action{
		Name: "StartRecording",
		Icon: "StartRecording",
		Command: "//recorder start",
	})

	acw.AddAction(widgets.Action{
		Name: "StopRecording",
		Icon: "StopRecording",
		Command: "//recorder stop",
	})

	return nil
}

func (p *RecorderPlugin) Unload() error {
	commandManager := app.GetCommandManager()

	commandManager.RemoveCommand("//recorder")
	commandManager.RemoveCommand("//export")
	commandManager.RemoveCommand("//recordings")

	acw := widgets.GetAdminControlsWidget()

	acw.RemoveAction("StartRecording")
	acw.RemoveAction("StopRecording")

	client := app.GetClient()

	for i, callback := range client.OnPlayerFinish {
		if callback.Key == "recording" {
			client.OnPlayerFinish = slices.Delete(client.OnPlayerFinish, i, i+1)
			break
		}
	}

	for i, callback := range client.OnPreEndRound {
		if callback.Key == "recording" {
			client.OnPreEndRound = slices.Delete(client.OnPreEndRound, i, i+1)
			break
		}
	}

	return nil
}

func (p *RecorderPlugin) startRecording(name string) {
	c := app.GetGoController()
	
	mode := c.MapManager.CurrentMode
	if mode == "Trackmania/TM_TimeAttack_Online.Script.txt" {
		c.Server.Client.OnPlayerFinish = append(c.Server.Client.OnPlayerFinish, gbxclient.GbxCallbackStruct[events.PlayerWayPointEventArgs]{
			Key:  "recording",
			Call: p.onPlayerFinish})
		p.createRecording(name, "TimeAttack")
	} else if mode == "Trackmania/TM_Rounds_Online.Script.txt" {
		c.Server.Client.OnPreEndRound = append(c.Server.Client.OnPreEndRound, gbxclient.GbxCallbackStruct[events.ScoresEventArgs]{
			Key:  "recording",
			Call: p.onPreEndRound})
		p.createRecording(name, "Rounds")
	} else {
		c.Server.Client.OnPreEndRound = append(c.Server.Client.OnPreEndRound, gbxclient.GbxCallbackStruct[events.ScoresEventArgs]{
			Key:  "recording",
			Call: p.onPreEndRoundMatch})
		p.createRecording(name, "Match")
	}

	p.IsRecording = true
}

func (p *RecorderPlugin) createRecording(name string, modeType string) {
	c := app.GetGoController()
	
	newRecording := database.NewRecording(database.Recording{
		Name: name,
		Mode: c.MapManager.CurrentMode,
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
	c := app.GetGoController()
	
	for i, callback := range c.Server.Client.OnPlayerFinish {
		if callback.Key == "recording" {
			c.Server.Client.OnPlayerFinish = slices.Delete(c.Server.Client.OnPlayerFinish, i, i+1)
		}
	}

	for i, callback := range c.Server.Client.OnPreEndRound {
		if callback.Key == "recording" {
			c.Server.Client.OnPreEndRound = slices.Delete(c.Server.Client.OnPreEndRound, i, i+1)
		}
	}

	p.IsRecording = false
	zap.L().Info("Recording stopped")
}

func (p *RecorderPlugin) onPlayerFinish(playerFinishEvent events.PlayerWayPointEventArgs) {
	c := app.GetGoController()
	
	mapId := c.MapManager.CurrentMapDB.ID

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
	c := app.GetGoController()
	
	mapId := c.MapManager.CurrentMapDB.ID

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
	c := app.GetGoController()
	
	mapId := c.MapManager.CurrentMapDB.ID

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
	c := app.GetGoController()
	
	if len(args) < 1 {
		go c.ChatUsage("//recorder [*start | stop] [name]", login)
		return
	}

	switch args[0] {
	case "start":
		if p.IsRecording {
			go c.ChatError("Already recording", nil, login)
			return
		}

		name := time.Now().Format("2006-01-02 15:04:05")
		if len(args) > 1 {
			name = strings.Join(args[1:], " ")
		}

		p.startRecording(name)
		go c.Chat("#Primary#Recording started", login)
	case "stop":
		if !p.IsRecording {
			go c.ChatError("Not recording", nil, login)
			return
		}

		p.stopRecording()
		go c.Chat("#Primary#Recording stopped with id #White#"+p.Recording.ID.Hex(), login)
	default:
		go c.ChatUsage("//recorder [*start | stop] [name]", login)
	}
}

func (p *RecorderPlugin) recordingsCommand(login string, args []string) {
	c := app.GetGoController()
	
	recordingsDB, err := database.GetRecordings(context.Background())
	if err != nil {
		zap.L().Error("Failed to get recordings", zap.Error(err))
		go c.ChatError("Failed to get recordings", nil, login)
		return
	}

	window := ui.NewGridWindow(&login)
	window.SetTemplate("recorder/recording.jet")
	window.Title = "Recordings"
	window.Items = make([]any, 0, len(recordingsDB))

	for _, recordingDB := range recordingsDB {
		var recording models.Recording
		recordingDB.ToModel(&recording)
		window.Actions[recording.ID] = app.GetUIManager().AddAction(p.handleDownloadAnswer, recording.ID)
		window.Items = append(window.Items, recording)
	}
	
	go window.Display()
}

func (p *RecorderPlugin) exportToCSVCommand(login string, args []string) {
	c := app.GetGoController()
	
	if len(args) < 1 {
		go c.ChatUsage("//export [*recording id]", login)
		return
	}

	recordingID := args[0]
	objectID, err := primitive.ObjectIDFromHex(recordingID)
	if err != nil {
		zap.L().Error("Invalid recording ID", zap.Error(err))
		go c.ChatError("Invalid recording ID", nil, login)
		return
	}

	err = p.exportToCSV(objectID)
	if err != nil {
		go c.ChatError("Failed to export recording to CSV", nil, login)
		return
	}

	go c.Chat("#Primary#Recording exported to CSV", login)
}

func (p *RecorderPlugin) handleDownloadAnswer(login string, data any, _ any) {
	c := app.GetGoController()
	
	if id, err := primitive.ObjectIDFromHex(data.(string)); err != nil {
		zap.L().Error("Invalid recording ID", zap.Error(err))
		go c.ChatError("Invalid recording ID", nil, login)
	} else {
		if err = p.exportToCSV(id); err != nil {
			go c.ChatError("Failed to export recording to CSV", nil, login)
		} else {
			go c.Chat("#Primary#Recording exported to CSV", login)
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
