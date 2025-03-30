package app

import (
	"sync"

	"github.com/MRegterschot/GbxRemoteGo/events"
	"github.com/MRegterschot/GbxRemoteGo/gbxclient"
	"github.com/MRegterschot/GbxRemoteGo/structs"
	"github.com/MRegterschot/GoController/database"
	"github.com/MRegterschot/GoController/utils"
	"go.uber.org/zap"
)

type MapManager struct {
	Maps         []structs.TMMapInfo
	CurrentMap   structs.TMMapInfo
	CurrentMapDB database.Map
	NextMap      structs.TMMapInfo
	PreviousMap  *structs.TMMapInfo
	CurrentMode  string
	MapsPath     string
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

	c := GetClient()

	mm.SyncMaps()

	c.OnBeginMap = append(c.OnBeginMap, gbxclient.GbxCallbackStruct[events.MapEventArgs]{
		Key:  "mapmanager",
		Call: mm.onBeginMap})
	c.OnMapListModified = append(c.OnMapListModified, gbxclient.GbxCallbackStruct[events.MapListModifiedEventArgs]{
		Key:  "mapmanager",
		Call: mm.onMapListModified})

	mode, err := c.GetScriptName()
	if err != nil {
		zap.L().Error("Failed to get mode script text", zap.Error(err))
		return
	}

	mm.CurrentMode = mode.CurrentValue
	mm.CurrentMap = mm.GetCurrentMapInfo()
	mm.NextMap = mm.GetNextMapInfo()

	mm.CurrentMapDB = GetDatabaseManager().SyncMap(mm.CurrentMap)

	mm.MapsPath, err = c.GetMapsDirectory()
	if err != nil {
		zap.L().Error("Failed to get maps directory", zap.Error(err))
	}

	zap.L().Info("MapManager initialized")
}

func (mm *MapManager) SyncMaps() {
	maps, err := GetClient().GetMapList(-1, 0)
	if err != nil {
		zap.L().Error("Failed to get map list", zap.Error(err))
		return
	}

	chunckedMaps := utils.ChunkArray(maps, 100)
	mapList := make([]structs.TMMapInfo, 0, len(maps))
	for _, chunk := range chunckedMaps {
		for _, m := range chunk {
			mapInfo, err := GetClient().GetMapInfo(m.FileName)
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

func (mm *MapManager) onBeginMap(mapEvent events.MapEventArgs) {
	mm.PreviousMap = &structs.TMMapInfo{
		UId:            mm.CurrentMap.UId,
		Name:           mm.CurrentMap.Name,
		FileName:       mm.CurrentMap.FileName,
		Author:         mm.CurrentMap.Author,
		AuthorNickname: mm.CurrentMap.AuthorNickname,
		Environnement:  mm.CurrentMap.Environnement,
		Mood:           mm.CurrentMap.Mood,
		BronzeTime:     mm.CurrentMap.BronzeTime,
		SilverTime:     mm.CurrentMap.SilverTime,
		GoldTime:       mm.CurrentMap.GoldTime,
		AuthorTime:     mm.CurrentMap.AuthorTime,
		CopperPrice:    mm.CurrentMap.CopperPrice,
		LapRace:        mm.CurrentMap.LapRace,
		NbLaps:         mm.CurrentMap.NbLaps,
		NbCheckpoints:  mm.CurrentMap.NbCheckpoints,
		MapType:        mm.CurrentMap.MapType,
		MapStyle:       mm.CurrentMap.MapStyle,
	}

	mm.CurrentMap = structs.TMMapInfo{
		UId:            mapEvent.Map.Uid,
		Name:           mapEvent.Map.Name,
		FileName:       mapEvent.Map.FileName,
		Author:         mapEvent.Map.Author,
		AuthorNickname: mapEvent.Map.AuthorNickname,
		Environnement:  mapEvent.Map.Environnement,
		Mood:           mapEvent.Map.Mood,
		BronzeTime:     mapEvent.Map.BronzeTime,
		SilverTime:     mapEvent.Map.SilverTime,
		GoldTime:       mapEvent.Map.GoldTime,
		AuthorTime:     mapEvent.Map.AuthorTime,
		CopperPrice:    mapEvent.Map.CopperPrice,
		LapRace:        mapEvent.Map.LapRace,
		NbLaps:         mapEvent.Map.NbLaps,
		NbCheckpoints:  mapEvent.Map.NbCheckpoints,
		MapType:        mapEvent.Map.MapType,
		MapStyle:       mapEvent.Map.MapStyle,
	}

	for i, m := range mm.Maps {
		if m.UId == mm.CurrentMap.UId {
			if i < len(mm.Maps)-1 {
				mm.NextMap = mm.Maps[i+1]
			} else {
				mm.NextMap = mm.Maps[0]
			}
			break
		}
	}

	mm.CurrentMapDB = GetDatabaseManager().SyncMap(mm.CurrentMap)

	if mode, err := GetClient().GetScriptName(); err != nil {
		zap.L().Error("Failed to get mode script text", zap.Error(err))
	} else {
		mm.CurrentMode = mode.CurrentValue
	}
}

func (mm *MapManager) onMapListModified(mapListModifiedEvent events.MapListModifiedEventArgs) {
	if mapListModifiedEvent.IsListModified {
		mm.SyncMaps()
	}
}

func (mm *MapManager) GetCurrentMapInfo() structs.TMMapInfo {
	currentMap, err := GetClient().GetCurrentMapInfo()
	if err != nil {
		zap.L().Error("Failed to get current map info", zap.Error(err))
		return structs.TMMapInfo{}
	}

	return currentMap
}

func (mm *MapManager) GetNextMapInfo() structs.TMMapInfo {
	nextMap, err := GetClient().GetNextMapInfo()
	if err != nil {
		zap.L().Error("Failed to get next map info", zap.Error(err))
		return structs.TMMapInfo{}
	}

	return nextMap
}
