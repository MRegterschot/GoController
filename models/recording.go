package models

import "time"

type Recording struct {
	ID   string
	Name string
	Type string
	Mode string
	Maps []MapRecords

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type MapRecords struct {
	ID          string
	Map         Map
	MatchRounds []MatchRound
	Rounds      []Round
	Finishes    []PlayerFinish
}

type MatchRound struct {
	ID          string
	RoundNumber int
	Teams       []Team
}

type Round struct {
	ID          string
	RoundNumber int
	Players     []PlayerRound
}

type Team struct {
	ID          string
	TeamID      int
	Name        string
	Points      int
	TotalPoints int
	Players     []PlayerRound
}

type PlayerRound struct {
	ID          string
	Player      Player
	Login       string
	AccountId   string
	Points      int
	TotalPoints int
	Time        int
	Checkpoints []int
}

type PlayerFinish struct {
	ID          string
	Player      Player
	Login       string
	AccountId   string
	Time        int
	Checkpoints []int
	Timestamp   time.Time
}
