package ui

type ListItem struct {
	Name        string
	Description string
	Value       string
}

type ListWindow struct {
	Window
}

func NewListWindow(login *string) *ListWindow {
	w := NewWindow(login)
	w.Template = "list.jet"

	lw := &ListWindow{
		Window: *w,
	}

	return lw
}
