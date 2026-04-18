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
	limit  int64
	offset int64
}

func NewPaginationParams(limit string, offset string) PaginationParams {
	return PaginationParams{
		limit:  cap(parseIntOr(limit, DefaultLimit), MaxLimit),
		offset: parseIntOr(offset, DefaultOffset),
	}
}

func (p PaginationParams) Limit() int64 {
	return p.limit
}

func (p PaginationParams) Offset() int64 {
	return p.offset
}

func cap(i int64, max int64) int64 {
	if i > max {
		return max
	}
	return i
}

func parseIntOr(s string, defaultValue int64) int64 {
	if s == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}

	return int64(i)
}
