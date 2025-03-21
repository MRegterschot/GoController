package models

import "time"

type Map struct {
	ID  string
	Name string
	UId string
	FileName string
	Author string
	AuthorNickname string
	AuthorTime int
	GoldTime int
	SilverTime int
	BronzeTime int

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}