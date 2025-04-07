package app

import (
	"context"
	"sync"

	"github.com/MRegterschot/GbxRemoteGo/structs"
	"github.com/MRegterschot/GoController/database"
	"github.com/MRegterschot/GoController/models"
	"github.com/MRegterschot/GoController/utils"
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
		var ubiUid string
		webIdentities, err := api.GetNadeoAPI().GetWebIdentities([]string{utils.DecodeSlug(player.Login)})
		if err != nil {
			zap.L().Error("Failed to get web identities", zap.Error(err))
		} else {
			if len(webIdentities) > 0 {
				ubiUid = webIdentities[0].Uid
			}
		}

		database.InsertPlayer(ctx, database.NewPlayer(database.Player{
			Login:    player.Login,
			NickName: player.NickName,
			Path:     player.Path,
			Roles:    []string{},
			UbiUid:   ubiUid,
		}))
	} else {
		dbPlayer.Update(database.Player{
			Login:    player.Login,
			NickName: player.NickName,
			Path:     player.Path,
			Roles:    dbPlayer.Roles,
			UbiUid:   dbPlayer.UbiUid,
		})

		database.UpdatePlayer(ctx, dbPlayer)
	}
}

func (dbm *DatabaseManager) SyncPlayers() {
	players := GetPlayerManager().Players
	logins := make([]string, 0)
	for _, player := range players {
		logins = append(logins, player.Login)
	}

	playersDB, err := database.GetPlayersByLogins(context.Background(), logins)
	if err != nil {
		zap.L().Error("Failed to get players from database", zap.Error(err))
		return
	}

	newAccountIds := make([]string, 0)
	newPlayers := make([]models.DetailedPlayer, 0)
	existingPlayers := make([]database.Player, 0)
	for _, player := range players {
		found := false
		for _, playerDB := range playersDB {
			if playerDB.Login == player.Login {
				found = true
				existingPlayers = append(existingPlayers, playerDB)
				if playerDB.UbiUid == "" {
					newAccountIds = append(newAccountIds, utils.DecodeSlug(player.Login))
				}
				break
			}
		}

		if !found {
			newPlayers = append(newPlayers, player)
			newAccountIds = append(newAccountIds, utils.DecodeSlug(player.Login))
		}
	}

	nadeoApi := api.GetNadeoAPI()
	webIdentities, err := nadeoApi.GetWebIdentities(newAccountIds)
	if err != nil {
		zap.L().Error("Failed to get web identities", zap.Error(err))
		return
	}

	for _, playerDB := range existingPlayers {
		for _, webIdentity := range webIdentities {
			if webIdentity.AccountId == utils.DecodeSlug(playerDB.Login) {
				playerDB.UbiUid = webIdentity.Uid
				break
			}
		}

		for _, player := range players {
			if player.Login == playerDB.Login {
				playerDB.Update(database.Player{
					Login:    player.Login,
					NickName: player.NickName,
					Path:     player.Path,
					Roles:    playerDB.Roles,
					UbiUid:   playerDB.UbiUid,
				})
				database.UpdatePlayer(context.Background(), playerDB)
				break
			}
		}
	}

	newPlayersDB := make([]database.Player, 0)
	for _, player := range newPlayers {
		var ubiUid string

		for _, webIdentity := range webIdentities {
			if webIdentity.AccountId == utils.DecodeSlug(player.Login) {
				ubiUid = webIdentity.Uid
				break
			}
		}

		newPlayersDB = append(newPlayersDB, database.NewPlayer(database.Player{
			Login:    player.Login,
			NickName: player.NickName,
			Path:     player.Path,
			Roles:    []string{},
			UbiUid:   ubiUid,
		}))
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
	} else if mapDB.ThumbnailUrl == "" {
		mapsInfo, err := api.GetNadeoAPI().GetMapsInfo([]string{mapInfo.UId})
		if err != nil {
			zap.L().Error("Failed to get maps info from Nadeo", zap.Error(err))
		}

		if len(mapsInfo) > 0 {
			mapDB.Submitter = mapsInfo[0].Submitter
			mapDB.FileUrl = mapsInfo[0].FileUrl
			mapDB.ThumbnailUrl = mapsInfo[0].ThumbnailUrl
			mapDB.Timestamp = primitive.NewDateTimeFromTime(mapsInfo[0].Timestamp)
			mapDB.Update(mapDB)
			database.UpdateMap(ctx, mapDB)
		}

		return mapDB
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
				break
			}
		}

		if !found {
			newUids = append(newUids, uid)
		}
	}

	if len(newUids) == 0 {
		zap.L().Info("No new maps to sync")
		return
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
