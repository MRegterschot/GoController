package models

import (
	"time"

	"github.com/MRegterschot/GbxRemoteGo/structs"
)

type Player struct {
	ID string
	Login string
	NickName string
	Path string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type DetailedPlayer struct {
	structs.TMPlayerDetailedInfo
	IsAdmin bool
}