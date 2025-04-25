package ui

import "github.com/MRegterschot/GoController/app"

type Widget struct {
	*app.Manialink
	Hidden bool
}

func NewWidget(template string) *Widget {
	ml := app.NewManialink(nil)
	ml.Template = template

	w := &Widget{
		Manialink: ml,
	}

	return w
}

func (w *Widget) Hide() {
	w.Hidden = true
	w.Manialink.Hide()
}

func (w *Widget) Display() {
	w.Hidden = false
	w.Manialink.Display()
}
