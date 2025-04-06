package app

import (
	"context"
	"sync"

	"github.com/MRegterschot/GbxRemoteGo/structs"
	"github.com/MRegterschot/GoController/database"
	"github.com/MRegterschot/GoController/models"
	"github.com/MRegterschot/GoController/utils/api"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type DatabaseManager struct {
}

var (
	dbmInstance *DatabaseManager
	dbmOnce     sync.Once
)

func GetDatabaseManager() *DatabaseManager {
	dbmOnce.Do(func() {
		dbmInstance = &DatabaseManager{}
	})
	return dbmInstance
}

func (dbm *DatabaseManager) Init() {
	zap.L().Info("Initializing DatabaseManager")
	err := database.Connect()
	if err != nil {
		zap.L().Fatal("Failed to connect to database", zap.Error(err))
	}
	zap.L().Info("Connected to database")
	zap.L().Info("DatabaseManager initialized")
}

func (dbm *DatabaseManager) SyncPlayer(player models.DetailedPlayer) {
	ctx := context.Background()
	dbPlayer, err := database.GetPlayerByLogin(ctx, player.Login)
	if err != nil {
		database.InsertPlayer(ctx, database.NewPlayer(database.Player{
			Login:    player.Login,
			NickName: player.NickName,
			Path:     player.Path,
			Roles:    []string{},
		}))
	} else {
		dbPlayer.Update(database.Player{
			Login:    player.Login,
			NickName: player.NickName,
			Path:     player.Path,
			Roles:    dbPlayer.Roles,
		})

		database.UpdatePlayer(ctx, dbPlayer)
	}
}

func (dbm *DatabaseManager) SyncPlayers() {
	players := GetPlayerManager().Players
	for _, player := range players {
		dbm.SyncPlayer(player)
	}
}

func (dbm *DatabaseManager) SyncMap(mapInfo structs.TMMapInfo) database.Map {
	ctx := context.Background()
	if mapDB, err := database.GetMapByUId(ctx, mapInfo.UId); err != nil {
		mapsInfo, err := api.GetNadeoAPI().GetMapsInfo([]string{mapInfo.UId})
		if err != nil {
			zap.L().Error("Failed to get maps info from Nadeo", zap.Error(err))
		}

		m := database.Map{
			Name:           mapInfo.Name,
			UId:            mapInfo.UId,
			FileName:       mapInfo.FileName,
			Author:         mapInfo.Author,
			AuthorNickname: mapInfo.AuthorNickname,
			AuthorTime:     mapInfo.AuthorTime,
			GoldTime:       mapInfo.GoldTime,
			SilverTime:     mapInfo.SilverTime,
			BronzeTime:     mapInfo.BronzeTime,
		}

		if len(mapsInfo) > 0 {
			m.Submitter = mapsInfo[0].Submitter
			m.FileUrl = mapsInfo[0].FileUrl
			m.ThumbnailUrl = mapsInfo[0].ThumbnailUrl
			m.Timestamp = primitive.NewDateTimeFromTime(mapsInfo[0].Timestamp)
		}

		newMap := database.NewMap(m)
		database.InsertMap(ctx, newMap)
		return newMap
	} else {
		return mapDB
	}
}

func (dbm *DatabaseManager) SyncMaps() {
	maps := GetMapManager().Maps
	uids := make([]string, 0)
	for _, mapInfo := range maps {
		uids = append(uids, mapInfo.UId)
	}

	mapsDB, err := database.GetMapsByUIds(context.Background(), uids)
	if err != nil {
		zap.L().Error("Failed to get maps from database", zap.Error(err))
		return
	}

	newUids := make([]string, 0)
	for _, uid := range uids {
		found := false
		for _, mapDB := range mapsDB {
			if mapDB.UId == uid {
				found = true
			}
		}

		if !found {
			newUids = append(newUids, uid)
		}
	}

	nadeoApi := api.GetNadeoAPI()

	mapsInfo, err := nadeoApi.GetMapsInfo(newUids)
	if err != nil {
		zap.L().Error("Failed to get maps info from Nadeo", zap.Error(err))
		return
	}

	newMaps := make([]database.Map, 0)
	for _, mapInfo := range mapsInfo {
		for _, m := range maps {
			if m.UId == mapInfo.MapUid {
				newMaps = append(newMaps, database.NewMap(database.Map{
					Name:           m.Name,
					UId:            m.UId,
					FileName:       m.FileName,
					Author:         m.Author,
					AuthorNickname: m.AuthorNickname,
					AuthorTime:     m.AuthorTime,
					GoldTime:       m.GoldTime,
					SilverTime:     m.SilverTime,
					BronzeTime:     m.BronzeTime,
					Submitter:      mapInfo.Submitter,
					Timestamp:      primitive.NewDateTimeFromTime(mapInfo.Timestamp),
					FileUrl:        mapInfo.FileUrl,
					ThumbnailUrl:   mapInfo.ThumbnailUrl,
				}))
				break
			}
		}
	}

	if len(newMaps) > 0 {
		if _, err := database.InsertMaps(context.Background(), newMaps); err != nil {
			zap.L().Error("Failed to insert maps into database", zap.Error(err))
		}
	}
}
