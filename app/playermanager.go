package app

import (
	"fmt"
	"regexp"
	"sync"

	"github.com/MRegterschot/GbxRemoteGo/events"
	"github.com/MRegterschot/GbxRemoteGo/gbxclient"
	"github.com/MRegterschot/GoController/models"
	"go.uber.org/zap"
	"slices"
)

var fakeplayerRe = regexp.MustCompile(`^\*fakeplayer\d+\*$`)

type PlayerManager struct {
	Players []models.Player
}

var (
	plmInstance *PlayerManager
	plmOnce     sync.Once
)

func GetPlayerManager() *PlayerManager {
	plmOnce.Do(func() {
		plmInstance = &PlayerManager{}
	})
	return plmInstance
}

func (plm *PlayerManager) Init() {
	zap.L().Info("Initializing PlayerManager")
	plm.SyncPlayers()
	GetClient().OnPlayerConnect = append(GetClient().OnPlayerConnect, gbxclient.GbxCallbackStruct[events.PlayerConnectEventArgs]{
		Key:  "pmPlayerConnect",
		Call: plm.onPlayerConnect})
	GetClient().OnPlayerDisconnect = append(GetClient().OnPlayerDisconnect, gbxclient.GbxCallbackStruct[events.PlayerDisconnectEventArgs]{
		Key:  "pmPlayerDisconnect",
		Call: plm.onPlayerDisconnect})
	zap.L().Info("PlayerManager initialized")
}

func (plm *PlayerManager) SyncPlayers() {
	players, err := GetClient().GetPlayerList(-1, 0)
	if err != nil {
		zap.L().Error("Failed to get player list", zap.Error(err))
		return
	}

	for _, player := range players {
		if player.PlayerId == 0 {
			continue
		}
		detailedInfo, err := GetClient().GetDetailedPlayerInfo(player.Login)
		if err != nil {
			zap.L().Error("Failed to get detailed player info", zap.Error(err))
			continue
		}
		plm.Players = append(plm.Players, models.Player{
			TMPlayerDetailedInfo: detailedInfo,
			IsAdmin:              GetGoController().IsAdmin(player.Login),
		})
	}

	GetDatabaseManager().SyncPlayers()
}

func (plm *PlayerManager) GetPlayer(login string) *models.Player {
	for i := range plm.Players {
		if plm.Players[i].Login == login {
			return &plm.Players[i] // Return the actual struct reference
		}
	}

	detailedInfo, err := GetClient().GetDetailedPlayerInfo(login)
	if err != nil {
		return nil
	}
	player := models.Player{
		TMPlayerDetailedInfo: detailedInfo,
		IsAdmin:              GetGoController().IsAdmin(login),
	}

	GetDatabaseManager().SyncPlayer(player)

	plm.Players = append(plm.Players, player)
	return &plm.Players[len(plm.Players)-1]
}

func (plm *PlayerManager) onPlayerConnect(playerConnectEvent events.PlayerConnectEventArgs) {
	if fakeplayerRe.MatchString(playerConnectEvent.Login) {
		return
	}

	for _, player := range plm.Players {
		if player.Login == playerConnectEvent.Login {
			go GetClient().Kick(player.Login, "You are already connected")
			return
		}
	}

	player := plm.GetPlayer(playerConnectEvent.Login)

	go GetGoController().Chat(fmt.Sprintf("Welcome %s!", player.NickName))
}

func (plm *PlayerManager) onPlayerDisconnect(playerDisconnectEvent events.PlayerDisconnectEventArgs) {
	for i, player := range plm.Players {
		if player.Login == playerDisconnectEvent.Login {
			plm.Players = slices.Delete(plm.Players, i, i+1)
			go GetGoController().Chat(fmt.Sprintf("%s disconnected", player.NickName))
			return
		}
	}
}
