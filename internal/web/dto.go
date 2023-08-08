package web

type CollectionDto[T any] struct {
	Total int `json:"total"`
	From  int `json:"from"`
	To    int `json:"to"`
	Items []T `json:"items"`
}

func NewCollectionDto[T any](items []T, from int, to int) *CollectionDto[T] {
	return &CollectionDto[T]{
		Items: items,
		Total: len(items),
		From:  from,
		To:    to,
	}
}
