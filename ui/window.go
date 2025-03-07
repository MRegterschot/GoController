package ui

import "github.com/MRegterschot/GoController/app"

type Window struct {
	app.Manialink
}

func NewWindow(login *string) *Window {
	ml := app.NewManialink(login)
	ml.Template = "window.jet"
	w := &Window{
		Manialink: *ml,
	}

	w.Actions["close"] = app.GetUIManager().AddAction(w.Hide, "")

	return w
}

func (w *Window) Hide(_ string, _ interface{}, _ interface{}) {
	w.Manialink.Destroy()
}