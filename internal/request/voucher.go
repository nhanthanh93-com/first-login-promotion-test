package request

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"trinity/helpers"
	"trinity/internal/model"
)

type CVoucherReq struct {
	Code string `json:"code" form:"code" validate:"required"`
}

func (r *CVoucherReq) Bind(c *gin.Context) (*CVoucherReq, error) {
	if err := c.ShouldBind(&r); err != nil {
		return nil, err
	}
	if err := helpers.Validate(r); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *CVoucherReq) Model() model.Voucher {
	var m model.Voucher
	m.Code = r.Code
	return m
}

type GVoucherReq struct {
	ID string `json:"id" form:"id" validate:"required"`
}

func (r *GVoucherReq) Bind(c *gin.Context) (*GVoucherReq, error) {
	id := c.Param("id")
	r.ID = id
	return r, nil
}

type LVoucherReq struct {
	ID   string `json:"id" form:"id"`
	Code string `json:"code" form:"code"`
	SortByDateRange
	Paginate
}

func (r *LVoucherReq) Bind(c *gin.Context) (*LVoucherReq, error) {
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

type UVoucherReq struct {
	ID     string `json:"id" form:"id" validate:"required"`
	IsUsed *bool  `json:"is_used,omitempty" form:"is_used" validate:"required"`
}

func (r *UVoucherReq) Bind(c *gin.Context) (*UVoucherReq, error) {
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

func (r *UVoucherReq) Model() model.Voucher {
	var m model.Voucher
	m.ID = uuid.MustParse(r.ID)
	m.IsUsed = r.IsUsed
	return m
}
