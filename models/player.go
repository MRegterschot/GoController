package models

import (
	"github.com/MRegterschot/GbxRemoteGo/structs"
)

type Player struct {
	structs.TMPlayerDetailedInfo
	IsAdmin bool
}