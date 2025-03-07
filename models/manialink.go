package models

import (
	"github.com/google/uuid"
)

type ManialinkAction struct {
	Callback func()
	Data     interface{}
}

type UISize struct {
	Width  int
	Height int
}

type UIPos struct {
	X int
	Y int
	Z int
}

type Manialink struct {
	ID           string
	Size         UISize
	Pos          UIPos
	Template     string
	Actions      map[string]ManialinkAction
	Data         interface{}
	Recipient    *string
	Title        string
	FirstDisplay bool
}

func NewManialink(login *string) *Manialink {
	return &Manialink{
		ID: uuid.NewString(),
		Size: UISize{
			Width:  150,
			Height: 120,
		},
		Pos: UIPos{
			X: 0,
			Y: 20,
			Z: 1,
		},
		Template:  "manialink.jet",
		Actions:   make(map[string]ManialinkAction),
		Data:      nil,
		Recipient: login,
		Title:     "",
		FirstDisplay: true,
	}
}

func (ml *Manialink) Display() {
	if ml.FirstDisplay {
		ml.FirstDisplay = false
	} else {
		
	}
}