package models

import "time"

type MapInfo struct {
	Author                   string    `json:"author"`
	AuthorScore              int       `json:"authorScore"`
	BronzeScore              int       `json:"bronzeScore"`
	CollectionName           string    `json:"collectionName"`
	CreatedWithGamepadEditor bool      `json:"createdWithGamepadEditor"`
	CreatedWithSimpleEditor  bool      `json:"createdWithSimpleEditor"`
	Filename                 string    `json:"filename"`
	GoldScore                int       `json:"goldScore"`
	IsPlayable               bool      `json:"isPlayable"`
	MapId                    string    `json:"mapId"`
	MapStyle                 string    `json:"mapStyle"`
	MapType                  string    `json:"mapType"`
	MapUid                   string    `json:"mapUid"`
	Name                     string    `json:"name"`
	SilverScore              int       `json:"silverScore"`
	Submitter                string    `json:"submitter"`
	Timestamp                time.Time `json:"timestamp"`
	FileUrl                  string    `json:"fileUrl"`
	ThumbnailUrl             string    `json:"thumbnailUrl"`
}

type WebIdentity struct {
	AccountId string    `json:"accountId"`
	Provider  string    `json:"provider"`
	Uid       string    `json:"uid"`
	Timestamp time.Time `json:"timestamp"`
}
