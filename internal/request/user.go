package request

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"trinity/helpers"
	"trinity/internal/model"
)

type RUserReq struct {
	Email      string `json:"email,omitempty" form:"email" validate:"required"`
	CampaignID string `json:"campaign_id,omitempty" form:"campaign_id" validate:"required"`
}

func (r *RUserReq) Bind(c *gin.Context) (*RUserReq, error) {
	if err := c.ShouldBind(&r); err != nil {
		return nil, err
	}
	if err := helpers.Validate(r); err != nil {
		return nil, err
	}

	ok := helpers.ValidateUUID(r.CampaignID)
	if !ok {
		return nil, fmt.Errorf("error: Invalid request parameters")
	}
	return r, nil
}

type CUserReq struct {
	Email string `json:"email" validate:"required"`
}

func (r *CUserReq) Bind(c *gin.Context) (*CUserReq, error) {
	if err := c.ShouldBind(&r); err != nil {
		return nil, err
	}
	if err := helpers.Validate(r); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *CUserReq) Model() model.User {
	var m model.User
	m.Email = r.Email
	return m
}

type GUserReq struct {
	ID string `json:"id,omitempty" validate:"required"`
}

func (r *GUserReq) Bind(c *gin.Context) (*GUserReq, error) {
	id := c.Param("id")
	r.ID = id
	return r, nil
}

type LUserReq struct {
	ID       string `json:"id,omitempty" form:"id"`
	Email    string `json:"email,omitempty" form:"email"`
	Vouchers string `json:"vouchers,omitempty" form:"vouchers"`
	SortByDateRange
	Paginate
}

func (r *LUserReq) Bind(c *gin.Context) (*LUserReq, error) {
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

type UUserReq struct {
	ID    string `json:"-,omitempty" form:"id" validate:"required"`
	Email string `json:"email,omitempty" form:"email" validate:"required"`
}

func (r *UUserReq) Bind(c *gin.Context) (*UUserReq, error) {
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

func (r *UUserReq) Model() model.User {
	var m model.User
	m.ID = uuid.MustParse(r.ID)
	m.Email = r.Email
	return m
}
