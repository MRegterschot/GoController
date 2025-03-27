package ui

import (
	"fmt"

	"github.com/MRegterschot/GoController/app"
	"github.com/MRegterschot/GoController/models"
)

type Grid struct {
	Cols int
	Rows int
	Gap  int
}

type GridWindow struct {
	*Window
	Grid
	Items      []any
	Pagination models.PaginationResult[any]
	Template   string
	AddData    func()
}

func NewGridWindow(login *string) *GridWindow {
	w := NewWindow(login)
	w.Template = "grid.jet"

	gw := &GridWindow{
		Window: w,
		Items:  []any{},
		Pagination: models.PaginationResult[any]{
			Items:       []any{},
			TotalItems:  0,
			CurrentPage: 0,
			TotalPages:  0,
			PageSize:    12,
		},
		Grid: Grid{Cols: 4, Rows: 3, Gap: 2},
		AddData: func() {},
	}

	uim := app.GetUIManager()
	gw.Actions["start"] = uim.AddAction(gw.paginate, "start")
	gw.Actions["previous"] = uim.AddAction(gw.paginate, "previous")
	gw.Actions["next"] = uim.AddAction(gw.paginate, "next")
	gw.Actions["end"] = uim.AddAction(gw.paginate, "end")

	return gw
}

func (gw *GridWindow) SetTemplate(template string) {
	gw.Window.SetTemplate(template)
}

func (gw *GridWindow) paginate(_ string, data any, _ any) {
	action, ok := data.(string)
	if !ok {
		return
	}

	if len(gw.Items) == 0 {
		return
	}

	gw.Pagination.UpdatePage(action)
	gw.Pagination.Paginate(gw.Items, gw.Pagination.CurrentPage, gw.Grid.Cols*gw.Grid.Rows)

	gw.Data = map[string]any{
		"Pagination": gw.Pagination,
		"Grid":       gw.Grid,
	}
}

func (gw *GridWindow) Display() {
	gw.paginate("", "start", nil)
	gw.AddData()
	fmt.Println(gw.Data)
	gw.Window.Display()
}
