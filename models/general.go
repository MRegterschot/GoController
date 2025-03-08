package models

type PaginationResult[T any] struct {
	Items       []T
	TotalItems  int
	CurrentPage int
	TotalPages  int
	PageSize    int
}
