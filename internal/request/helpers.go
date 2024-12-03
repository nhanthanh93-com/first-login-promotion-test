package request

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math"
	"time"
)

const (
	DefaultLimit  int32  = 10
	DefaultSortBy string = "created_at"
)

type Paginate struct {
	Page     int32  `json:"-" form:"page"`
	Limit    int32  `json:"-" form:"limit"`
	SortBy   string `json:"-" form:"sort_type"`
	SortType string `json:"-" form:"sort_by"`
}

func (r *Paginate) GetSortType() string {
	switch r.SortType {
	case "desc":
		return "desc"
	case "asc":
		return "asc"
	default:
		return "asc"
	}
}

func (r *Paginate) GetTotalPage(count int64) int32 {
	limit := r.GetLimit()
	return int32(math.Ceil(float64(count) / float64(limit)))
}

func (r *Paginate) GetSortBy() string {
	if r.SortBy == "" {
		return DefaultSortBy
	}
	return r.SortBy
}

func (r *Paginate) GetOrder() string {
	return fmt.Sprintf("%v %v", r.GetSortBy(), r.GetSortType())
}

func (r *Paginate) GetOffset() int32 {
	return (r.GetPage() - 1) * r.GetLimit()
}

func (r *Paginate) GetPage() int32 {
	if r.Page <= 0 {
		return 1
	}
	return r.Page
}

func (r *Paginate) GetLimit() int32 {
	if r.Limit <= 0 {
		return DefaultLimit
	}
	return r.Limit
}

func (r *Paginate) GetSkip() int32 {
	return (r.GetPage() - 1) * r.GetLimit()
}
func (r *Paginate) FormatPagination(totalItems, perPage int64) gin.H {
	totalPage := r.GetTotalPage(totalItems)

	return gin.H{
		"current_page": r.Page,
		"per_page":     perPage,
		"total_page":   totalPage,
		"total_item":   totalItems,
	}
}

type SortByDateRange struct {
	StarDateStr string     `json:"-" form:"start_date"`
	EndDateStr  string     `json:"-" form:"end_date"`
	StartDate   *time.Time `json:"start_date,omitempty"`
	EndDate     *time.Time `json:"end_date,omitempty"`
}

func (r *SortByDateRange) BuildSortByDateRange() (*SortByDateRange, error) {
	start, err := time.Parse("2006-01-02", r.StarDateStr)
	if err != nil {
		return nil, err
	}
	end, err := time.Parse("2006-01-02", r.EndDateStr)
	if err != nil {
		return nil, err
	}
	return &SortByDateRange{
		StartDate: &start,
		EndDate:   &end,
	}, nil
}
