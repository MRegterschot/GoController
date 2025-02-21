package app

import (
	"sync"

	"github.com/MRegterschot/GbxRemoteGo/events"
	"github.com/MRegterschot/GbxRemoteGo/gbxclient"
	"github.com/MRegterschot/GbxRemoteGo/structs"
	"github.com/MRegterschot/GoController/utils"
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
	GetGoController().Server.Client.OnBeginMap = append(GetGoController().Server.Client.OnBeginMap, mm.onBeginMap)
	zap.L().Info("MapManager initialized")
}

func (mm *MapManager) SyncMaps() {
	maps, err := GetGoController().Server.Client.GetMapList(-1, 0)
	if err != nil {
		zap.L().Error("Failed to get map list", zap.Error(err))
		return
	}

	chunckedMaps := utils.ChunkArray(maps, 100)
	mapList := make([]structs.TMMapInfo, 0)
	for _, chunk := range chunckedMaps {
		for _, m := range chunk {
			mapInfo, err := GetGoController().Server.Client.GetMapInfo(m.FileName)
			if err != nil {
				zap.L().Error("Failed to get map info", zap.Error(err))
				continue
			}
			mapList = append(mapList, mapInfo)
		}
	}

	mm.Maps = mapList
	GetDatabaseManager().SyncMaps()
}

func (mm *MapManager) AddMap(mapInfo structs.TMMapInfo) {
	for _, m := range mm.Maps {
		if m.UId == mapInfo.UId {
			return
		}
	}
	mm.Maps = append(mm.Maps, mapInfo)
}

func (mm *MapManager) RemoveMap(uid string) {
	for i, m := range mm.Maps {
		if m.UId == uid {
			mm.Maps = append(mm.Maps[:i], mm.Maps[i+1:]...)
			return
		}
	}
}

func (mm *MapManager) GetMap(uid string) *structs.TMMapInfo {
	for i := range mm.Maps {
		if mm.Maps[i].UId == uid {
			return &mm.Maps[i]
		}
	}
	return nil
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

func (mm *MapManager) onMapListModified(_ *gbxclient.GbxClient, mapListModifiedEvent events.MapListModifiedEventArgs) {
	if mapListModifiedEvent.IsListModified {
		mm.SyncMaps()
	}
}