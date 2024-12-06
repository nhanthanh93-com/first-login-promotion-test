package handler

import (
	"github.com/gin-gonic/gin"
	"trinity/api/service"
	"trinity/helpers"
	"trinity/helpers/response"
	"trinity/internal/model"
	"trinity/internal/request"
)

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(productService service.ProductService) ProductHandler {
	return ProductHandler{productService}
}

// @Summary Create a new Product
// @Description Creates a new Product with the provided Product information.
// @Tags products
// @Accept json
// @Produce json
// @Param user body request.CProductReq true "Product information"
// @Success 200 {object} response.Response[model.Product]
// @Failure 500 {object} response.Response[any]
// @Failure 400 {object} response.Response[any]
// @Failure 404 {object} response.Response[any]
// @Router /products/create [post]
func (h *ProductHandler) Create(c *gin.Context) {
	var req request.CProductReq
	r, err := req.Bind(c)
	if err != nil {
		response.InvalidError(c, "Invalid request parameters")
		return
	}
	m := r.Model()

	product, err := h.productService.Create(m)
	if err != nil {
		response.HandleGormError(c, err)
		return
	}

	response.SuccessResponse[*model.Product](c, "Created", product)
}

// @Summary Get Product by ID
// @Description Get the details of a Product based on the Product ID.
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID" example("67ae81a6-5284-436a-a2e5-54c3ebeaa241")
// @Success 200 {object} response.Response[model.Product]
// @Failure 500 {object} response.Response[any]
// @Failure 400 {object} response.Response[any]
// @Failure 404 {object} response.Response[any]
// @Router /products/{id} [get]
func (h *ProductHandler) Find(c *gin.Context) {
	var req request.GProductReq
	r, err := req.Bind(c)
	if err != nil {
		response.InvalidError(c, "Invalid request parameters")
		return
	}

	product, err := h.productService.Find(r)
	if err != nil {
		response.NotFoundError(c, "Resource not found")
		return
	}

	response.SuccessResponse[*model.Product](c, "Success", product)
}

// @Summary Get list all Products
// @Description API Get list all Products
// @Tags products
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit of users per page" default(10)
// @Param name query string false "name"
// @Param start_date query string false "Start date" format(date-time) example("2024-12-04T00:00:00Z")
// @Param end_date query string false "End date" format(date-time) example("2024-12-04T00:00:00Z")
// @Success 200 {object} response.Response[response.ResponseData[[]model.Product]]
// @Failure 500 {object} response.Response[any]
// @Failure 400 {object} response.Response[any]
// @Failure 404 {object} response.Response[any]
// @Router /products/list [get]
func (h *ProductHandler) List(c *gin.Context) {
	var req request.LProductReq
	r, err := req.Bind(c)
	if err != nil {
		response.InvalidError(c, "Invalid request parameters")
		return
	}

	wg := helpers.NewWgGroup()
	var rs []*model.Product
	var count int64

	wg.Go(func() error {
		rs, err = h.productService.List(r)
		if err == nil {
			return err
		}
		return nil
	})

	wg.Go(func() error {
		count, err = h.productService.Count(r)
		if err == nil {
			return err
		}
		return nil
	})

	if err = wg.Wait(); err != nil {
		response.HandleGormError(c, err)
		return
	}

	data := response.ResponseData[[]*model.Product]{
		Pagination: response.PaginationResponse(req.GetPage(), req.GetLimit(), int32(len(rs)), count),
		Results:    &rs,
	}

	response.SuccessResponse[response.ResponseData[[]*model.Product]](c, "List", data)

}

// @Summary Update Product information
// @Description Update Product details based on the Product ID.
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID" example("67ae81a6-5284-436a-a2e5-54c3ebeaa241")
// @Param Product body request.UProductReq true "Product information"
// @Success 200 {object} response.Response[model.Product]
// @Failure 500 {object} response.Response[any]
// @Failure 400 {object} response.Response[any]
// @Failure 404 {object} response.Response[any]
// @Router /products/{id} [put]
func (h *ProductHandler) Update(c *gin.Context) {
	var req request.UProductReq
	r, err := req.Bind(c)
	if err != nil {
		response.InvalidError(c, "Invalid request parameters")
		return
	}
	m := r.Model()

	Product, err := h.productService.Update(m)
	if err != nil {
		response.HandleGormError(c, err)
		return
	}

	response.SuccessResponse[*model.Product](c, "Updated", Product)
}

// @Summary Delete Product information
// @Description Delete Product details based on the Product ID.
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID" example("67ae81a6-5284-436a-a2e5-54c3ebeaa241")
// @Success 200 {object} response.Response[model.Product]
// @Failure 500 {object} response.Response[any]
// @Failure 400 {object} response.Response[any]
// @Failure 404 {object} response.Response[any]
// @Router /products/{id} [delete]
func (h *ProductHandler) Delete(c *gin.Context) {
	var req request.GProductReq
	r, err := req.Bind(c)
	if err != nil {
		response.InvalidError(c, "Invalid request parameters")
		return
	}

	ok, err := h.productService.Delete(r)
	if err != nil || !ok {
		response.HandleGormError(c, err)
		return

	}

	response.SuccessResponse[interface{}](c, "Deleted", nil)
}
