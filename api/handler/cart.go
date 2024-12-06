package handler

import (
	"github.com/gin-gonic/gin"
	"trinity/api/service"
	"trinity/helpers/response"
	"trinity/internal/model"
	"trinity/internal/request"
)

type CartHandler struct {
	cartService service.CartService
}

func NewCartHandler(cartService service.CartService) CartHandler {
	return CartHandler{cartService}
}

// @Summary Add to cart
// @Description Add new Cart Item.
// @Tags carts
// @Accept json
// @Produce json
// @Param cart body request.AddToCartReq true "Product information"
// @Success 200 {object} response.Response[model.Order]
// @Failure 500 {object} response.Response[any]
// @Failure 400 {object} response.Response[any]
// @Failure 404 {object} response.Response[any]
// @Router /carts/add [post]
func (h *CartHandler) AddToCart(c *gin.Context) {
	var req request.AddToCartReq
	r, err := req.Bind(c)
	if err != nil {
		response.InvalidError(c, "Invalid request parameters")
		return
	}
	m := r.Model()

	cart, err := h.cartService.AddToCart(m)
	if err != nil {
		response.HandleGormError(c, err)
	}

	response.SuccessResponse[*model.CartItem](c, "Added", cart)
}

// @Summary Get Cart by ID
// @Description Get the details of a Cart based on the Cart ID.
// @Tags carts
// @Accept json
// @Produce json
// @Param id path string true "Cart ID" example("67ae81a6-5284-436a-a2e5-54c3ebeaa241")
// @Success 200 {object} response.Response[model.Cart]
// @Failure 500 {object} response.Response[any]
// @Failure 400 {object} response.Response[any]
// @Failure 404 {object} response.Response[any]
// @Router /carts/{id} [get]
func (h *CartHandler) Find(c *gin.Context) {
	var req request.GCartReq
	r, err := req.Bind(c)
	if err != nil {
		response.InvalidError(c, "Invalid request parameters")
		return
	}

	cart, err := h.cartService.Find(r)
	if err != nil {
		response.HandleGormError(c, err)
	}

	response.SuccessResponse[*model.Cart](c, "Success", cart)
}

// @Summary Create a new Order
// @Description Creates a new Order with the provided Order information.
// @Tags orders
// @Accept json
// @Produce json
// @Param user_id path string true "User ID" example("67ae81a6-5284-436a-a2e5-54c3ebeaa241")
// @Param user body request.COrderReq true "Product information"
// @Success 200 {object} response.Response[model.Order]
// @Failure 500 {object} response.Response[any]
// @Failure 400 {object} response.Response[any]
// @Failure 404 {object} response.Response[any]
// @Router /orders/{user_id} [post]
func (h *CartHandler) CreateOrder(c *gin.Context) {
	var req request.COrderReq
	r, err := req.Bind(c)
	if err != nil {
		response.InvalidError(c, "Invalid request parameters")
		return
	}
	order, err := h.cartService.CreateOrder(r)
	if err != nil {
		response.HandleGormError(c, err)
	}

	response.SuccessResponse[*model.Order](c, "Success", order)
}

// @Summary Update Order information
// @Description Update Order details based on the Order ID.
// @Tags orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID" example("67ae81a6-5284-436a-a2e5-54c3ebeaa241")
// @Param user body request.UOrderStatusReq true "Product information"
// @Success 200 {object} response.Response[model.Order]
// @Failure 500 {object} response.Response[any]
// @Failure 400 {object} response.Response[any]
// @Failure 404 {object} response.Response[any]
// @Router /orders/{id}/status [put]
func (h *CartHandler) UpdateOrderStatus(c *gin.Context) {
	var req request.UOrderStatusReq
	r, err := req.Bind(c)
	if err != nil {
		response.InvalidError(c, "Invalid request parameters")
		return
	}
	order, err := h.cartService.UpdateOrderStatus(r)
	if err != nil {
		response.HandleGormError(c, err)
	}

	response.SuccessResponse[*model.Order](c, "Updated", order)

}

// @Summary Delete a cart item
// @Description Deletes a cart item based on the provided cart item ID.
// @Tags carts
// @Accept json
// @Produce json
// @Param id path string true "Cart Item ID" example("67ae81a6-5284-436a-a2e5-54c3ebeaa241")
// @Success 200 {object} response.Response[any]
// @Failure 500 {object} response.Response[any]
// @Failure 400 {object} response.Response[any]
// @Failure 404 {object} response.Response[any]
// @Router /carts/{id} [delete]
func (h *CartHandler) DeleteCartItem(c *gin.Context) {
	var req request.GCartItemReq
	r, err := req.Bind(c)
	if err != nil {
		response.InvalidError(c, "Invalid request parameters")
		return
	}

	ok, err := h.cartService.DeleteCartItem(r)
	if err != nil || !ok {
		response.HandleGormError(c, err)
	}

	response.SuccessResponse[interface{}](c, "Deleted", nil)
}
