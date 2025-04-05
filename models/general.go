package models

type Toggle struct {
	Label  string
	Color  string
	Action string
}

type PaginationResult[T any] struct {
	Items       []T
	TotalItems  int
	CurrentPage int
	TotalPages  int
	PageSize    int
}

// Paginate an array
func (pr *PaginationResult[T]) Paginate(array []T, page int, pageSize int) {
	currPage := min(page+1, (len(array)+pageSize-1)/pageSize) - 1
	start := max(currPage*pageSize, 0)
	end := start + pageSize
	if start > len(array) {
		start = len(array)
	}
	if end > len(array) {
		end = len(array)
	}

	pr.Items = array[start:end]
	pr.TotalItems = len(array)
	pr.CurrentPage = currPage
	pr.TotalPages = (len(array) + pageSize - 1) / pageSize
	pr.PageSize = pageSize
}

// Update current page based on action
func (pr *PaginationResult[T]) UpdatePage(action string) {
	switch action {
	case "start":
		pr.CurrentPage = 0
	case "previous":
		if pr.CurrentPage > 0 {
			pr.CurrentPage--
		}
	case "next":
		if pr.CurrentPage < pr.TotalPages-1 {
			pr.CurrentPage++
		}
	case "end":
		pr.CurrentPage = pr.TotalPages - 1
	}
}

type Styling map[string]string
type Fonts map[string]string
type Icons map[string]string

type Theme struct {
	Styling Styling
	Fonts   Fonts
	Icons   Icons
}
