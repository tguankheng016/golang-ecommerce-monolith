package pagination

type PageResultDto[T any] struct {
	TotalCount int `json:"totalCount,omitempty" bson:"totalCount"`
	Items      []T `json:"items,omitempty" bson:"items"`
}

func NewPageResultDto[T any](items []T, totalCount int) PageResultDto[T] {
	pageResultDto := PageResultDto[T]{Items: items, TotalCount: totalCount}

	return pageResultDto
}
