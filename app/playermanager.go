package app

import (
	"sync"

	"github.com/MRegterschot/GbxRemoteGo/events"
	"github.com/MRegterschot/GbxRemoteGo/gbxclient"
	"github.com/MRegterschot/GoController/models"
	"go.uber.org/zap"
)

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
	GetGoController().Server.Client.OnPlayerConnect = append(GetGoController().Server.Client.OnPlayerConnect, plm.onPlayerConnect)
	GetGoController().Server.Client.OnPlayerDisconnect = append(GetGoController().Server.Client.OnPlayerDisconnect, plm.onPlayerDisconnect)
	zap.L().Info("PlayerManager initialized")
}

func (plm *PlayerManager) GetPlayer(login string) *models.Player {
	for _, p := range plm.Players {
		if p.Login == login {
			return &p
		}
	}

	detailedInfo, err := GetGoController().Server.Client.GetDetailedPlayerInfo(login)
	if err != nil {
		return nil
	}
	player := models.Player{
		TMPlayerDetailedInfo: detailedInfo,
		IsAdmin:              GetGoController().IsAdmin(login),
	}
	plm.Players = append(plm.Players, player)
	return &player
}

func (plm *PlayerManager) onPlayerConnect(client *gbxclient.GbxClient, playerConnectEvent events.PlayerConnectEventArgs) {
	for _, player := range plm.Players {
		if player.Login == playerConnectEvent.Login {
			go client.Kick(player.Login, "You are already connected")
			return
		}
	}

	go GetGoController().Chat("Welcome to the server!", playerConnectEvent.Login)
}

func (plm *PlayerManager) onPlayerDisconnect(_ *gbxclient.GbxClient, playerDisconnectEvent events.PlayerDisconnectEventArgs) {
	for i, player := range plm.Players {
		if player.Login == playerDisconnectEvent.Login {
			plm.Players = append(plm.Players[:i], plm.Players[i+1:]...)
			return
		}
	}
}