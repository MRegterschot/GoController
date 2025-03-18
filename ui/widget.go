package ui

import "github.com/MRegterschot/GoController/app"

type Widget struct {
	*app.Manialink
}

func NewWidget(template string) *Widget {
	ml := app.NewManialink(nil)
	ml.Template = template

	w := &Widget{
		Manialink: ml,
	}

	return w
}