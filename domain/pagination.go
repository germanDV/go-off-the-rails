package domain

import (
	"strconv"
)

const (
	DefaultLimit  = 100
	DefaultOffset = 0
	MaxLimit      = 999
)

type PaginationParams struct {
	limit  int
	offset int
}

func NewPaginationParams(limit string, offset string) PaginationParams {
	return PaginationParams{
		limit:  cap(parseIntOr(limit, DefaultLimit), MaxLimit),
		offset: parseIntOr(offset, DefaultOffset),
	}
}

func (p PaginationParams) Limit() int {
	return p.limit
}

func (p PaginationParams) Offset() int {
	return p.offset
}

func cap(i int, max int) int {
	if i > max {
		return max
	}
	return i
}

func parseIntOr(s string, defaultValue int) int {
	if s == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}

	return i
}
