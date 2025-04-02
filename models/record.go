package models

import "time"

type Record struct {
	ID     string
	Login  string
	Time   int
	MapUId string
	Player DetailedPlayer

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
