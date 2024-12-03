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

	response.SuccessResponse(c, "Success", voucher)
}

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

	response.SuccessResponse(c, "List", gin.H{
		"results":    rs,
		"pagination": r.FormatPagination(count, int64(len(rs))),
	})

}

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

	response.SuccessResponse(c, "Updated", voucher)
}

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

	response.SuccessResponse(c, "Deleted", nil)
}
