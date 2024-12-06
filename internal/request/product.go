package request

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"trinity/helpers"
	"trinity/internal/model"
)

type CProductReq struct {
	Name  string  `json:"name,omitempty" form:"name" validate:"required"`
	Price float64 `json:"price,omitempty" form:"price" validate:"required"`
	Stock int     `json:"stock,omitempty" form:"stock" validate:"required"`
}

func (r *CProductReq) Bind(c *gin.Context) (*CProductReq, error) {
	if err := c.ShouldBind(&r); err != nil {
		return nil, err
	}
	if err := helpers.Validate(r); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *CProductReq) Model() model.Product {
	var m model.Product
	m.Name = r.Name
	m.Price = r.Price
	return m
}

type GProductReq struct {
	ID string `json:"id,omitempty" form:"id" validate:"required"`
}

func (r *GProductReq) Bind(c *gin.Context) (*GProductReq, error) {
	id := c.Param("id")
	r.ID = id
	return r, nil
}

type LProductReq struct {
	ID    string `json:"id,omitempty" form:"id"`
	Name  string `json:"name,omitempty" form:"name"`
	Price *int32 `json:"price,omitempty" form:"price"`
	SortByDateRange
	Paginate
}

func (r *LProductReq) Bind(c *gin.Context) (*LProductReq, error) {
	if err := c.ShouldBindQuery(&r); err != nil {
		return nil, err
	}
	if err := helpers.Validate(r); err != nil {
		return nil, err
	}

	if r.StarDateStr != "" && r.EndDateStr != "" {
		dateRange, err := r.BuildSortByDateRange()
		if err != nil {
			return nil, err
		}
		r.StartDate = dateRange.StartDate
		r.EndDate = dateRange.EndDate
	}

	r.Limit = r.GetLimit()
	r.Page = r.GetPage()
	r.SortType = r.GetSortType()
	r.SortBy = r.GetSortBy()
	return r, nil
}

type UProductReq struct {
	ID    string  `json:"-,omitempty" form:"id" validate:"required"`
	Name  string  `json:"name,omitempty" form:"name" validate:"required"`
	Price float64 `json:"price,omitempty" form:"price" validate:"required"`
}

func (r *UProductReq) Bind(c *gin.Context) (*UProductReq, error) {
	id := c.Param("id")
	r.ID = id
	if err := c.ShouldBind(&r); err != nil {
		return nil, err
	}
	if err := helpers.Validate(r); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *UProductReq) Model() model.Product {
	var m model.Product
	m.ID = uuid.MustParse(r.ID)
	m.Name = r.Name
	m.Price = r.Price
	return m
}
