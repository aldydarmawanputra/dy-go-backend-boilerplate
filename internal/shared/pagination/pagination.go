package pagination

import "strconv"

const (
	defaultLimit = 20
	maxLimit     = 100
)

type Params struct {
	Limit  int
	Offset int
}

func Parse(limitStr, offsetStr string) Params {
	limit := atoiOr(limitStr, defaultLimit)
	if limit <= 0 || limit > maxLimit {
		limit = defaultLimit
	}
	offset := atoiOr(offsetStr, 0)
	if offset < 0 {
		offset = 0
	}
	return Params{Limit: limit, Offset: offset}
}

type Meta struct {
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
	Count  int   `json:"count"`
	Total  int64 `json:"total"`
}

func NewMeta(p Params, count int, total int64) Meta {
	return Meta{Limit: p.Limit, Offset: p.Offset, Count: count, Total: total}
}

func atoiOr(s string, fallback int) int {
	if n, err := strconv.Atoi(s); err == nil {
		return n
	}
	return fallback
}
