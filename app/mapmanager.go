package app

import (
	"sync"

	"github.com/MRegterschot/GbxRemoteGo/events"
	"github.com/MRegterschot/GbxRemoteGo/gbxclient"
	"github.com/MRegterschot/GbxRemoteGo/structs"
	"go.uber.org/zap"
)

type MapManager struct {
	Maps       []structs.TMMapInfo
	CurrentMap structs.TMMapInfo
	NextMap    structs.TMMapInfo
}

var (
	mmInstance *MapManager
	mmOnce     sync.Once
)

func GetMapManager() *MapManager {
	mmOnce.Do(func() {
		mmInstance = &MapManager{}
	})
	return mmInstance
}

func (mm *MapManager) Init() {
	zap.L().Info("Initializing MapManager")
	mm.SyncMaps()
	zap.L().Info("MapManager initialized")
}

func (mm *MapManager) SyncMaps() {
	maps, err := GetGoController().Server.Client.GetMapList(-1, 0)
	if err != nil {
		zap.L().Error("Failed to get map list", zap.Error(err))
		return
	}

	mm.Maps = maps
}

func (mm *MapManager) onBeginMap(_ *gbxclient.GbxClient, mapEvent events.MapEventArgs) {
	mm.CurrentMap = structs.TMMapInfo{
		UId: mapEvent.Map.Uid,
		Name: mapEvent.Map.Name,
		FileName: mapEvent.Map.FileName,
		Author: mapEvent.Map.Author,
		AuthorNickname: mapEvent.Map.AuthorNickname,
		Environnement: mapEvent.Map.Environnement,
		Mood: mapEvent.Map.Mood,
		BronzeTime: mapEvent.Map.BronzeTime,
		SilverTime: mapEvent.Map.SilverTime,
		GoldTime: mapEvent.Map.GoldTime,
		AuthorTime: mapEvent.Map.AuthorTime,
		CopperPrice: mapEvent.Map.CopperPrice,
		LapRace: mapEvent.Map.LapRace,
		NbLaps: mapEvent.Map.NbLaps,
		NbCheckpoints: mapEvent.Map.NbCheckpoints,
		MapType: mapEvent.Map.MapType,
		MapStyle: mapEvent.Map.MapStyle,
	}
	
	for i, m := range mm.Maps {
		if m.UId == mm.CurrentMap.UId {
			if i < len(mm.Maps) - 1 {
				mm.NextMap = mm.Maps[i + 1]
			} else {
				mm.NextMap = mm.Maps[0]
			}
			break
		}
	}
}
