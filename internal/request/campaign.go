package request

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"time"
	"trinity/helpers"
	"trinity/internal/model"
)

type CCampaignReq struct {
	Name      string    `json:"name,omitempty" form:"name" validate:"required"`
	MaxUser   int32     `json:"max_user,omitempty" form:"max_user" validate:"required"`
	ExpiresAt time.Time `json:"expires_at,omitempty" form:"expires_at" validate:"required"`
	Discount  int32     `json:"discount,omitempty" form:"discount" validate:"required"`
}

func (r *CCampaignReq) Bind(c *gin.Context) (*CCampaignReq, error) {
	if err := c.ShouldBind(&r); err != nil {
		return nil, err
	}
	if err := helpers.Validate(r); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *CCampaignReq) Model() model.Campaign {
	var m model.Campaign
	m.Name = r.Name
	m.MaxUser = r.MaxUser
	return m
}

type GCampaignReq struct {
	ID string `json:"id,omitempty" form:"id" validate:"required"`
}

func (r *GCampaignReq) Bind(c *gin.Context) (*GCampaignReq, error) {
	id := c.Param("id")
	r.ID = id
	return r, nil
}

type LCampaignReq struct {
	ID       string `json:"id,omitempty" form:"id"`
	Name     string `json:"name,omitempty" form:"name"`
	MaxUser  *int32 `json:"max_user,omitempty" form:"max_user"`
	Vouchers string `json:"vouchers,omitempty" form:"vouchers"`
	SortByDateRange
	Paginate
}

func (r *LCampaignReq) Bind(c *gin.Context) (*LCampaignReq, error) {
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

type UCampaignReq struct {
	ID      string `json:"-,omitempty" form:"id" validate:"required"`
	Name    string `json:"name,omitempty" form:"name" validate:"required"`
	MaxUser int32  `json:"max_user,omitempty" form:"max_user" validate:"required"`
}

func (r *UCampaignReq) Bind(c *gin.Context) (*UCampaignReq, error) {
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

func (r *UCampaignReq) Model() model.Campaign {
	var m model.Campaign
	m.ID = uuid.MustParse(r.ID)
	m.Name = r.Name
	m.MaxUser = r.MaxUser
	return m
}
