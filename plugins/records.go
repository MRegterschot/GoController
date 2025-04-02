package plugins

import (
	"context"
	"sync"

	"slices"

	"github.com/MRegterschot/GbxRemoteGo/events"
	"github.com/MRegterschot/GbxRemoteGo/gbxclient"
	"github.com/MRegterschot/GoController/app"
	"github.com/MRegterschot/GoController/database"
	"github.com/MRegterschot/GoController/models"
	"github.com/MRegterschot/GoController/plugins/widgets"
	"go.uber.org/zap"
)

type RecordsPlugin struct {
	Name         string
	Dependencies []string
	Loaded       bool

	LocalRecords  map[string]models.Record
	RecordsWidget *widgets.RecordsWidget
}

var (
	records     *RecordsPlugin
	recordsOnce sync.Once
)

func CreateRecordsPlugin() *RecordsPlugin {
	return &RecordsPlugin{
		Name:          "Records",
		Dependencies:  []string{},
		Loaded:        false,
		LocalRecords:  make(map[string]models.Record, 0),
		RecordsWidget: widgets.GetRecordsWidget(),
	}
}

func GetRecordsPlugin() *RecordsPlugin {
	recordsOnce.Do(func() {
		records = CreateRecordsPlugin()
	})
	return records
}

func (p *RecordsPlugin) Load() error {
	client := app.GetClient()

	client.OnPlayerFinish = append(client.OnPlayerFinish, gbxclient.GbxCallbackStruct[events.PlayerWayPointEventArgs]{
		Key:  "records",
		Call: p.onPlayerFinish,
	})

	client.OnBeginMap = append(client.OnBeginMap, gbxclient.GbxCallbackStruct[events.MapEventArgs]{
		Key:  "records",
		Call: p.onBeginMap,
	})

	currentMap, err := app.GetClient().GetCurrentMapInfo()
	if err != nil {
		zap.L().Error("Failed to get current map info", zap.Error(err))
		return err
	}

	mapUid := currentMap.UId
	records, err := database.GetRecordsByMapUId(context.Background(), mapUid)
	if err != nil {
		zap.L().Error("Failed to get records by map UId", zap.String("mapUId", mapUid), zap.Error(err))
	}

	for _, record := range records {
		if _, ok := p.LocalRecords[record.Login]; ok {
			// Record already exists, check if this record is better
			if p.LocalRecords[record.Login].Time > record.Time {
				p.updateLocalRecord(record)
			}
			continue
		}

		var localRecord models.Record
		record.ToModel(&localRecord)
		player := app.GetPlayerManager().GetPlayer(record.Login)
		if player != nil {
			localRecord.Player = *player
		}
		p.LocalRecords[record.Login] = localRecord
	}

	p.RecordsWidget.SetRecords(p.LocalRecords)

	return nil
}

func (p *RecordsPlugin) Unload() error {
	client := app.GetClient()

	for i, callback := range client.OnPlayerFinish {
		if callback.Key == "records" {
			client.OnPlayerFinish = slices.Delete(client.OnPlayerFinish, i, i+1)
			break
		}
	}

	for i, callback := range client.OnBeginMap {
		if callback.Key == "records" {
			client.OnBeginMap = slices.Delete(client.OnBeginMap, i, i+1)
			break
		}
	}

	return nil
}

func (p *RecordsPlugin) onPlayerFinish(playerFinishEvent events.PlayerWayPointEventArgs) {
	c := app.GetGoController()

	mapUid := c.MapManager.CurrentMap.UId

	record := database.NewRecord(database.Record{
		Login:  playerFinishEvent.Login,
		Time:   playerFinishEvent.RaceTime,
		MapUId: mapUid,
	})

	_, err := database.InsertRecord(context.Background(), record)
	if err != nil {
		zap.L().Error("Failed to insert record", zap.Any("record", record), zap.Error(err))
		return
	}

	if lr, ok := p.LocalRecords[record.Login]; ok {
		// Record already exists, check if this record is better
		if record.Time >= lr.Time {
			return
		}

		p.updateLocalRecord(record)
	} else {
		// Add new record to LocalRecords
		var newRecord models.Record
		record.ToModel(&newRecord)
		p.LocalRecords[record.Login] = newRecord
	}

	p.RecordsWidget.SetRecords(p.LocalRecords)
}

func (p *RecordsPlugin) onBeginMap(mapEvent events.MapEventArgs) {
	records, err := database.GetRecordsByMapUId(context.Background(), mapEvent.Map.Uid)
	if err != nil {
		zap.L().Error("Failed to get records by map UId", zap.String("mapUId", mapEvent.Map.Uid), zap.Error(err))
		return
	}

	for _, record := range records {
		if _, ok := p.LocalRecords[record.Login]; ok {
			// Record already exists, check if this record is better
			if p.LocalRecords[record.Login].Time > record.Time {
				p.updateLocalRecord(record)
			}
			continue
		}

		var localRecord models.Record
		record.ToModel(&localRecord)
		player := app.GetPlayerManager().GetPlayer(record.Login)
		if player != nil {
			localRecord.Player = *player
		}
		p.LocalRecords[record.Login] = localRecord
	}

	p.RecordsWidget.SetRecords(p.LocalRecords)
}

func (p *RecordsPlugin) updateLocalRecord(record database.Record) {
	newRecord := p.LocalRecords[record.Login]
	newRecord.Time = record.Time
	newRecord.CreatedAt = record.CreatedAt.Time()
	p.LocalRecords[record.Login] = newRecord
}

func init() {
	recordsPlugin := GetRecordsPlugin()
	app.GetPluginManager().PreLoadPlugin(recordsPlugin)
}
