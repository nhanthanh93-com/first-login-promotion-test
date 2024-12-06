package request

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"trinity/helpers"
	"trinity/internal/model"
)

type AddToCartReq struct {
	CartID    string `json:"cart_id,omitempty" form:"cart_id" validate:"required"`
	ProductID string `json:"product_id,omitempty" form:"product_id" validate:"required"`
	Quantity  int    `json:"quantity,omitempty" form:"quantity" validate:"required"`
}

func (r *AddToCartReq) Bind(c *gin.Context) (*AddToCartReq, error) {
	if err := c.ShouldBind(&r); err != nil {
		return nil, err
	}
	if err := helpers.Validate(r); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *AddToCartReq) Model() model.CartItem {
	var m model.CartItem
	m.CartID = uuid.MustParse(r.CartID)
	m.ProductID = uuid.MustParse(r.ProductID)
	return m
}

type GCartItemReq struct {
	ID string `json:"id,omitempty" validate:"required"`
}

func (r *GCartItemReq) Bind(c *gin.Context) (*GCartItemReq, error) {
	id := c.Param("id")
	r.ID = id
	return r, nil
}

type UCartItemReq struct {
	ID       string `json:"-,omitempty" form:"id" validate:"required"`
	Quantity int    `json:"quantity,omitempty" form:"quantity" validate:"required"`
}

func (r *UCartItemReq) Bind(c *gin.Context) (*UCartItemReq, error) {
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

func (r *UCartItemReq) Model() model.CartItem {
	var m model.CartItem
	m.ID = uuid.MustParse(r.ID)
	m.Quantity = r.Quantity
	return m
}

type GCartReq struct {
	ID string `json:"id,omitempty" validate:"required"`
}

func (r *GCartReq) Bind(c *gin.Context) (*GCartReq, error) {
	id := c.Param("id")
	r.ID = id
	return r, nil
}
