package request

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"trinity/helpers"
	"trinity/internal/model"
)

type COrderReq struct {
	CartID string  `json:"cart_id,omitempty" form:"cart_id" validate:"required"`
	UserID string  `json:"-,omitempty" form:"user_id" validate:"required"`
	Total  float64 `json:"total,omitempty" form:"total" validate:"required"`
	Status string  `json:"status,omitempty" form:"status" validate:"required"`
}

func (r *COrderReq) Bind(c *gin.Context) (*COrderReq, error) {
	userID := c.Param("user_id")
	r.UserID = userID
	if err := c.ShouldBind(&r); err != nil {
		return nil, err
	}
	if err := helpers.Validate(r); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *COrderReq) Model() model.Order {
	var m model.Order
	m.UserID = uuid.MustParse(r.UserID)
	m.Total = r.Total
	m.Status = r.Status
	return m
}

type UOrderStatusReq struct {
	ID     string `json:"-,omitempty" form:"id" validate:"required"`
	Status string `json:"status,omitempty" form:"status" validate:"required"`
}

func (r *UOrderStatusReq) Bind(c *gin.Context) (*UOrderStatusReq, error) {
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

func (r *UOrderStatusReq) Model() model.Order {
	var m model.Order
	m.ID = uuid.MustParse(r.ID)
	m.Status = r.Status
	return m
}
