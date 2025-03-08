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
	Actions      map[string]string
	Data         interface{}
	Recipient    *string
	Title        string
	FirstDisplay bool
}

func NewManialink(login *string) *Manialink {
	return &Manialink{
		ID: uuid.NewString(),
		Size: UISize{
			Width:  200,
			Height: 124,
		},
		Pos: UIPos{
			X: 0,
			Y: 20,
			Z: 1,
		},
		Template:  "manialink.jet",
		Actions:   make(map[string]string),
		Data:      nil,
		Recipient: login,
		Title:     "",
		FirstDisplay: true,
	}
}

func (ml *Manialink) Display() {
	if ml.FirstDisplay {
		ml.FirstDisplay = false
		GetUIManager().DisplayManialink(ml)
	} else {
		GetUIManager().RefreshManialink(ml)
	}
}

func (ml *Manialink) Hide() {
	GetUIManager().HideManialink(ml)
}

func (ml *Manialink) Destroy() {
	GetUIManager().DestroyManialink(ml)
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