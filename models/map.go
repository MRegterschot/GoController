package models

import (
	"time"

	"github.com/MRegterschot/GbxRemoteGo/structs"
)

type Map struct {
	ID             string
	Name           string
	UId            string
	FileName       string
	Author         string
	AuthorNickname string
	AuthorTime     int
	GoldTime       int
	SilverTime     int
	BronzeTime     int

	// Fields from nadeo API
	Submitter    string
	Timestamp    time.Time
	FileUrl      string
	ThumbnailUrl string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type QueueMap struct {
	Name             string
	UId              string
	FileName         string
	Author           string
	AuthorNickname   string
	QueuedBy         string
	QueuedByNickname string
	QueuedAt         time.Time
}

func (qm *QueueMap) ToQueueMap(m structs.TMMapInfo) {
	qm.Name = m.Name
	qm.UId = m.UId
	qm.FileName = m.FileName
	qm.Author = m.Author
	qm.AuthorNickname = m.AuthorNickname
	qm.QueuedAt = time.Now()
}
