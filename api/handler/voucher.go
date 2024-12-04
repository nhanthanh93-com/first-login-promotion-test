package handler

import (
	"github.com/gin-gonic/gin"
	"trinity/api/service"
	"trinity/helpers"
	"trinity/helpers/response"
	"trinity/internal/model"
	"trinity/internal/request"
)

type VoucherHandler struct {
	voucherService service.VoucherService
}

func NewVoucherHandler(voucherService service.VoucherService) VoucherHandler {
	return VoucherHandler{
		voucherService: voucherService,
	}
}

// @Summary Get voucher by ID
// @Description Get the details of a voucher based on the voucher ID.
// @Tags vouchers
// @Accept json
// @Produce json
// @Param id path string true "Voucher ID" example("67ae81a6-5284-436a-a2e5-54c3ebeaa241")
// @Success 200 {object} response.Response[model.Voucher]
// @Failure 500 {object} response.Response[any]
// @Failure 400 {object} response.Response[any]
// @Failure 404 {object} response.Response[any]
// @Router /vouchers/{id} [get]
func (h *VoucherHandler) Find(c *gin.Context) {
	var req request.GVoucherReq
	r, err := req.Bind(c)
	if err != nil {
		response.InvalidError(c, "Invalid request parameters")
		return
	}

	voucher, err := h.voucherService.Find(r)
	if err != nil {
		response.NotFoundError(c, "Resource not found")
		return
	}

	response.SuccessResponse[*model.Voucher](c, "Success", voucher)
}

// @Summary Get list all vouchers
// @Description API Get list all vouchers
// @Tags vouchers
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit of users per page" default(10)
// @Param code query string false "code"
// @Param start_date query string false "Start date" format(date-time) example("2024-12-04T00:00:00Z")
// @Param end_date query string false "End date" format(date-time) example("2024-12-04T00:00:00Z")
// @Success 200 {object} response.Response[response.ResponseData[[]model.Voucher]]
// @Failure 500 {object} response.Response[any]
// @Failure 400 {object} response.Response[any]
// @Failure 404 {object} response.Response[any]
// @Router /vouchers/list [get]
func (h *VoucherHandler) List(c *gin.Context) {
	var req request.LVoucherReq
	r, err := req.Bind(c)
	if err != nil {
		response.InvalidError(c, "Invalid request parameters")
		return
	}

	wg := helpers.NewWgGroup()
	var rs []*model.Voucher
	var count int64

	wg.Go(func() error {
		rs, err = h.voucherService.List(r)
		if err == nil {
			return err
		}
		return nil
	})

	wg.Go(func() error {
		count, err = h.voucherService.Count(r)
		if err == nil {
			return err
		}
		return nil
	})

	if err = wg.Wait(); err != nil {
		response.HandleGormError(c, err)
		return
	}

	data := response.ResponseData[[]*model.Voucher]{
		Pagination: response.PaginationResponse(req.GetPage(), req.GetLimit(), int32(len(rs)), count),
		Results:    &rs,
	}

	response.SuccessResponse[response.ResponseData[[]*model.Voucher]](c, "List", data)

}

// @Summary Update voucher information
// @Description Update voucher details based on the voucher ID.
// @Tags vouchers
// @Accept json
// @Produce json
// @Param id path string true "Voucher ID" example("67ae81a6-5284-436a-a2e5-54c3ebeaa241")
// @Param voucher body request.UVoucherReq true "Voucher information"
// @Success 200 {object} response.Response[model.Voucher]
// @Failure 500 {object} response.Response[any]
// @Failure 400 {object} response.Response[any]
// @Failure 404 {object} response.Response[any]
// @Router /vouchers/{id} [put]
func (h *VoucherHandler) Update(c *gin.Context) {
	var req request.UVoucherReq
	r, err := req.Bind(c)
	if err != nil {
		response.InvalidError(c, "Invalid request parameters")
		return
	}
	m := r.Model()

	voucher, err := h.voucherService.Update(m)
	if err != nil {
		response.HandleGormError(c, err)
		return
	}

	response.SuccessResponse[*model.Voucher](c, "Updated", voucher)
}

// @Summary Delete voucher information
// @Description Delete voucher details based on the campaign ID.
// @Tags vouchers
// @Accept json
// @Produce json
// @Param id path string true "Voucher ID" example("67ae81a6-5284-436a-a2e5-54c3ebeaa241")
// @Success 200 {object} response.Response[model.Voucher]
// @Failure 500 {object} response.Response[any]
// @Failure 400 {object} response.Response[any]
// @Failure 404 {object} response.Response[any]
// @Router /vouchers/{id} [delete]
func (h *VoucherHandler) Delete(c *gin.Context) {
	var req request.GVoucherReq
	r, err := req.Bind(c)
	if err != nil {
		response.InvalidError(c, "Invalid request parameters")
		return
	}

	ok, err := h.voucherService.Delete(r)
	if err != nil || !ok {
		response.HandleGormError(c, err)
		return

	}

	response.SuccessResponse[interface{}](c, "Deleted", nil)
}
