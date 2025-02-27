package app

import (
	"context"
	"sync"

	"github.com/MRegterschot/GbxRemoteGo/structs"
	"github.com/MRegterschot/GoController/database"
	"github.com/MRegterschot/GoController/models"
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

func (dbm *DatabaseManager) SyncPlayer(player models.Player) {
	ctx := context.Background()
	dbPlayer, err := database.GetPlayerByLogin(ctx, player.Login)
	if err != nil {
		database.InsertPlayer(ctx, database.NewPlayer(database.Player{
			Login:    player.Login,
			NickName: player.NickName,
			Path:     player.Path,
		}))
	} else {
		dbPlayer.Update(database.Player{
			Login:    player.Login,
			NickName: player.NickName,
			Path:     player.Path,
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
		newMap := database.NewMap(database.Map{
			Name:           mapInfo.Name,
			UId:            mapInfo.UId,
			FileName:       mapInfo.FileName,
			Author:         mapInfo.Author,
			AuthorNickname: mapInfo.AuthorNickname,
			AuthorTime:     mapInfo.AuthorTime,
			GoldTime:       mapInfo.GoldTime,
			SilverTime:     mapInfo.SilverTime,
			BronzeTime:     mapInfo.BronzeTime,
		})
		database.InsertMap(ctx, newMap)
		return newMap
	} else {
		return mapDB
	}
}

func (dbm *DatabaseManager) SyncMaps() {
	maps := GetMapManager().Maps
	for _, mapInfo := range maps {
		dbm.SyncMap(mapInfo)
	}
}
