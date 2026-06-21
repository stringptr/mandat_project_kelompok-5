package pagination

type Meta struct {
	CurrentPage int32 `json:"current_page"`
	PerPage     int32 `json:"per_page"`
	Total       int32 `json:"total"`
	LastPage    int32 `json:"last_page"`
}

func NewMeta(page, perPage, total int32) Meta {
	lastPage := total / perPage
	if total%perPage != 0 {
		lastPage++
	}
	if lastPage < 1 {
		lastPage = 1
	}
	return Meta{
		CurrentPage: page,
		PerPage:     perPage,
		Total:       total,
		LastPage:    lastPage,
	}
}

func ValidatePage(page int) int {
	if page < 1 {
		return 1
	}
	return page
}

func ValidatePerPage(perPage int) int {
	if perPage < 1 {
		return 20
	}
	if perPage > 100 {
		return 100
	}
	return perPage
}

func Offset(page, perPage int) int32 {
	return int32((page - 1) * perPage)
}
