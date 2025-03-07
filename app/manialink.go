package app

import (
	"bytes"

	"github.com/google/uuid"
)

type ManialinkAction struct {
	Callback func(string, interface{}, interface{})
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

func (ml *Manialink) Render() (string, error) {
	t, err := GetUIManager().Templates.GetTemplate(ml.Template)
	if err != nil {
		return "", err
	}

	data := map[string]interface{}{
		"ID": 	ml.ID,
		"Size": ml.Size,
		"Pos": 	ml.Pos,
		"Actions": ml.Actions,
		"Data": ml.Data,
		"Title": ml.Title,
		"Recipient": ml.Recipient,
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, nil, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}