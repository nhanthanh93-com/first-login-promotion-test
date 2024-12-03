package handler

import (
	"github.com/gin-gonic/gin"
	"trinity/api/service"
	"trinity/helpers"
	"trinity/helpers/response"
	"trinity/internal/model"
	"trinity/internal/request"
)

type CampaignHandler struct {
	campaignService service.CampaignService
}

func NewCampaignHandler(campaignService service.CampaignService) CampaignHandler {
	return CampaignHandler{
		campaignService: campaignService,
	}
}

// @Summary Create a new campaign
// @Description Creates a new campaign with the provided campaign information.
// @Tags campaigns
// @Accept json
// @Produce json
// @Param user body request.CCampaignReq true "Campaign information"
// @Success 200 {object} response.Response[model.Campaign]
// @Failure 500 {object} response.Response[any]
// @Failure 400 {object} response.Response[any]
// @Failure 404 {object} response.Response[any]
// @Router /campaigns/create [post]
func (h *CampaignHandler) Create(c *gin.Context) {
	var req request.CCampaignReq
	r, err := req.Bind(c)
	if err != nil {
		response.InvalidError(c, "Invalid request parameters")
		return
	}
	m := r.Model()

	campaign, err := h.campaignService.Create(m)
	if err != nil {
		response.HandleGormError(c, err)
		return
	}

	response.SuccessResponse[*model.Campaign](c, "Created", campaign)
}

// @Summary Get campaign by ID
// @Description Get the details of a campaign based on the campaign ID.
// @Tags campaigns
// @Accept json
// @Produce json
// @Param id path string true "Campaign ID" example("67ae81a6-5284-436a-a2e5-54c3ebeaa241")
// @Success 200 {object} response.Response[model.Campaign]
// @Failure 500 {object} response.Response[any]
// @Failure 400 {object} response.Response[any]
// @Failure 404 {object} response.Response[any]
// @Router /campaigns/{id} [get]
func (h *CampaignHandler) Find(c *gin.Context) {
	var req request.GCampaignReq
	r, err := req.Bind(c)
	if err != nil {
		response.InvalidError(c, "Invalid request parameters")
		return
	}

	campaign, err := h.campaignService.Find(r)
	if err != nil {
		response.NotFoundError(c, "Resource not found")
		return
	}

	response.SuccessResponse[*model.Campaign](c, "Success", campaign)
}

// @Summary Get list all campaigns
// @Description API Get list all campaigns
// @Tags campaigns
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit of users per page" default(10)
// @Param name query string false "name"
// @Param max_user query string false "Max User"
// @Param vouchers query string false "vouchers" example("Tkv5bqH9Y4fdd0CjhGvBo, 2KOf7dAhJVRBKCbmZ8J9i")
// @Param start_date query string false "Start date" format(date-time) example("2024-12-04T00:00:00Z")
// @Param end_date query string false "End date" format(date-time) example("2024-12-04T00:00:00Z")
// @Success 200 {object} response.Response[response.ResponseData[[]model.Campaign]]
// @Failure 500 {object} response.Response[any]
// @Failure 400 {object} response.Response[any]
// @Failure 404 {object} response.Response[any]
// @Router /users/list [get]
func (h *CampaignHandler) List(c *gin.Context) {
	var req request.LCampaignReq
	r, err := req.Bind(c)
	if err != nil {
		response.InvalidError(c, "Invalid request parameters")
		return
	}

	wg := helpers.NewWgGroup()
	var rs []*model.Campaign
	var count int64

	wg.Go(func() error {
		rs, err = h.campaignService.List(r)
		if err == nil {
			return err
		}
		return nil
	})

	wg.Go(func() error {
		count, err = h.campaignService.Count(r)
		if err == nil {
			return err
		}
		return nil
	})

	if err = wg.Wait(); err != nil {
		response.HandleGormError(c, err)
		return
	}

	data := response.ResponseData[[]*model.Campaign]{
		Pagination: response.PaginationResponse(req.GetPage(), req.GetLimit(), int32(len(rs)), count),
		Results:    &rs,
	}

	response.SuccessResponse[response.ResponseData[[]*model.Campaign]](c, "List", data)

}

// @Summary Update campaign information
// @Description Update campaign details based on the campaign ID.
// @Tags campaigns
// @Accept json
// @Produce json
// @Param id path string true "Campaign ID" example("67ae81a6-5284-436a-a2e5-54c3ebeaa241")
// @Param campaign body request.UCampaignReq true "Campaign information"
// @Success 200 {object} response.Response[model.Campaign]
// @Failure 500 {object} response.Response[any]
// @Failure 400 {object} response.Response[any]
// @Failure 404 {object} response.Response[any]
// @Router /campaigns/{id} [put]
func (h *CampaignHandler) Update(c *gin.Context) {
	var req request.UCampaignReq
	r, err := req.Bind(c)
	if err != nil {
		response.InvalidError(c, "Invalid request parameters")
		return
	}
	m := r.Model()

	campaign, err := h.campaignService.Update(m)
	if err != nil {
		response.HandleGormError(c, err)
		return
	}

	response.SuccessResponse[*model.Campaign](c, "Updated", campaign)
}

// @Summary Delete campaign information
// @Description Delete campaign details based on the campaign ID.
// @Tags campaigns
// @Accept json
// @Produce json
// @Param id path string true "Campaign ID" example("67ae81a6-5284-436a-a2e5-54c3ebeaa241")
// @Success 200 {object} response.Response[model.Campaign]
// @Failure 500 {object} response.Response[any]
// @Failure 400 {object} response.Response[any]
// @Failure 404 {object} response.Response[any]
// @Router /campaigns/{id} [delete]
func (h *CampaignHandler) Delete(c *gin.Context) {
	var req request.GCampaignReq
	r, err := req.Bind(c)
	if err != nil {
		response.InvalidError(c, "Invalid request parameters")
		return
	}

	ok, err := h.campaignService.Delete(r)
	if err != nil || !ok {
		response.HandleGormError(c, err)
		return

	}

	response.SuccessResponse[interface{}](c, "Deleted", nil)
}
