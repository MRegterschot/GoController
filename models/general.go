package models

type PaginationResult[T any] struct {
	Items       []T
	TotalItems  int
	CurrentPage int
	TotalPages  int
	PageSize    int
}

// Paginate an array
func (pr *PaginationResult[T]) Paginate(array []T, page int, pageSize int) {
	start := page * pageSize
	end := start + pageSize
	if start > len(array) {
		start = len(array)
	}
	if end > len(array) {
		end = len(array)
	}

	pr.Items = array[start:end]
	pr.TotalItems = len(array)
	pr.CurrentPage = page
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