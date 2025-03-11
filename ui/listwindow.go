package ui

import (
	"github.com/MRegterschot/GoController/app"
	"github.com/MRegterschot/GoController/models"
	"github.com/MRegterschot/GoController/utils"
)

type Column struct {
	Name  string
	Width int    // percentage of the total width
	Type  string // text, input
}

type ListWindow struct {
	*Window
	Columns     []Column
	Items       [][]any
	Pagination  models.PaginationResult[[]any]
	UpdateItems func([][]any, any)
}

func NewListWindow(login *string) *ListWindow {
	w := NewWindow(login)
	w.Template = "list.jet"

	lw := &ListWindow{
		Window: w,
		Items:  [][]any{},
		Pagination: models.PaginationResult[[]any]{
			Items:       [][]any{},
			TotalItems:  0,
			CurrentPage: 0,
			TotalPages:  0,
			PageSize:    14,
		},
	}

	uim := app.GetUIManager()
	lw.Actions["start"] = uim.AddAction(lw.paginate, "start")
	lw.Actions["previous"] = uim.AddAction(lw.paginate, "previous")
	lw.Actions["next"] = uim.AddAction(lw.paginate, "next")
	lw.Actions["end"] = uim.AddAction(lw.paginate, "end")

	return lw
}

func (lw *ListWindow) paginate(_ string, data any, entries any) {
	action, ok := data.(string)
	if !ok {
		return
	}

	if len(lw.Items) == 0 {
		return
	}

	switch action {
	case "start":
		lw.Pagination.CurrentPage = 0
	case "previous":
		if lw.Pagination.CurrentPage > 0 {
			lw.Pagination.CurrentPage--
		} else {
			return
		}
	case "next":
		if lw.Pagination.CurrentPage < lw.Pagination.TotalPages-1 {
			lw.Pagination.CurrentPage++
		} else {
			return
		}
	case "end":
		lw.Pagination.CurrentPage = lw.Pagination.TotalPages - 1
	}

	lw.UpdateItems(lw.Items, entries)
	lw.Pagination = utils.Paginate(lw.Items, lw.Pagination.CurrentPage, lw.Pagination.PageSize)

	lw.Data = struct {
		Columns    []Column
		Pagination models.PaginationResult[[]any]
	}{
		Columns:    lw.Columns,
		Pagination: lw.Pagination,
	}

	lw.Window.Display()
}

func (lw *ListWindow) Display() {
	lw.paginate("", "start", nil)
}
