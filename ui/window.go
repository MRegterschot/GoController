package ui

import (
	"github.com/MRegterschot/GoController/app"
)

type Window struct {
	*app.Manialink
}

func NewWindow(login *string) *Window {
	ml := app.NewManialink(login)
	ml.Template = "window.jet"
	w := &Window{
		Manialink: ml,
	}

	w.Actions["close"] = app.GetUIManager().AddAction(w.Destroy, nil)

	return w
}

func (w *Window) SetTemplate(template string) {
	w.Manialink.Template = template
}

func (w *Window) Destroy(_ string, _ any, _ any) {
	w.Manialink.Destroy()
}
