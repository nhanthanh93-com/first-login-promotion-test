package handler

import (
	"github.com/gin-gonic/gin"
	"trinity/api/service"
	"trinity/helpers"
	"trinity/helpers/response"
	"trinity/internal/model"
	"trinity/internal/request"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) UserHandler {
	return UserHandler{
		userService: userService,
	}
}

// @Summary Register a new user
// @Description Registers a new user with the provided registration details, such as email and campaign ID.
// @Tags promo
// @Accept json
// @Produce json
// @Param user body request.RUserReq true "User information"
// @Success 200 {object} response.Response[model.User]
// @Failure 500 {object} response.Response[any]
// @Failure 400 {object} response.Response[any]
// @Failure 404 {object} response.Response[any]
// @Router /promo/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req request.RUserReq
	r, err := req.Bind(c)
	if err != nil {
		response.InvalidError(c, "Invalid request parameters")
		return
	}
	user, err := h.userService.Register(r)
	if err != nil {
		response.HandleGormError(c, err)
		return
	}

	response.SuccessResponse[*model.User](c, "Registered", user)
}

// @Summary Create a new user
// @Description Creates a new user with the provided user information.
// @Tags users
// @Accept json
// @Produce json
// @Param user body request.CUserReq true "User information"
// @Success 200 {object} response.Response[model.User]
// @Failure 500 {object} response.Response[any]
// @Failure 400 {object} response.Response[any]
// @Failure 404 {object} response.Response[any]
// @Router /users/create [post]
func (h *UserHandler) Create(c *gin.Context) {
	var req request.CUserReq
	r, err := req.Bind(c)
	if err != nil {
		response.InvalidError(c, "Invalid request parameters")
		return
	}
	m := r.Model()

	user, err := h.userService.Create(m)
	if err != nil {
		response.HandleGormError(c, err)
		return
	}

	response.SuccessResponse[*model.User](c, "Created", user)
}

// @Summary Get user by ID
// @Description Get the details of a user based on the user ID.
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID" example("67ae81a6-5284-436a-a2e5-54c3ebeaa241")
// @Success 200 {object} response.Response[model.User]
// @Failure 500 {object} response.Response[any]
// @Failure 400 {object} response.Response[any]
// @Failure 404 {object} response.Response[any]
// @Router /users/{id} [get]
func (h *UserHandler) Find(c *gin.Context) {
	var req request.GUserReq
	r, err := req.Bind(c)
	if err != nil {
		response.InvalidError(c, "Invalid request parameters")
		return
	}

	user, err := h.userService.Find(r)
	if err != nil {
		response.NotFoundError(c, "Resource not found")
		return
	}

	response.SuccessResponse[*model.User](c, "Success", user)
}

// @Summary Get list all users
// @Description API Get list all users
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit of users per page" default(10)
// @Param email query string false "email" example("john.doe@example.com")
// @Param vouchers query string false "vouchers" example("Tkv5bqH9Y4fdd0CjhGvBo, 2KOf7dAhJVRBKCbmZ8J9i")
// @Param start_date query string false "Start date" format(date-time) example("2024-12-04T00:00:00Z")
// @Param end_date query string false "End date" format(date-time) example("2024-12-04T00:00:00Z")
// @Success 200 {object} response.Response[response.ResponseData[[]model.User]]
// @Failure 500 {object} response.Response[any]
// @Failure 400 {object} response.Response[any]
// @Failure 404 {object} response.Response[any]
// @Router /users/list [get]
func (h *UserHandler) List(c *gin.Context) {
	var req request.LUserReq
	r, err := req.Bind(c)
	if err != nil {
		response.InvalidError(c, "Invalid request parameters")
		return
	}

	wg := helpers.NewWgGroup()
	var rs []*model.User
	var count int64

	wg.Go(func() error {
		rs, err = h.userService.List(r)
		if err == nil {
			return err
		}
		return nil
	})

	wg.Go(func() error {
		count, err = h.userService.Count(r)
		if err == nil {
			return err
		}
		return nil
	})

	if err = wg.Wait(); err != nil {
		response.HandleGormError(c, err)
		return
	}

	data := response.ResponseData[[]*model.User]{
		Pagination: response.PaginationResponse(req.GetPage(), req.GetLimit(), int32(len(rs)), count),
		Results:    &rs,
	}

	response.SuccessResponse[response.ResponseData[[]*model.User]](c, "List user", data)
}

// @Summary Update user information
// @Description Update user details based on the user ID.
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID" example("67ae81a6-5284-436a-a2e5-54c3ebeaa241")
// @Param user body request.UUserReq true "User information"
// @Success 200 {object} response.Response[model.User]
// @Failure 500 {object} response.Response[any]
// @Failure 400 {object} response.Response[any]
// @Failure 404 {object} response.Response[any]
// @Router /users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	var req request.UUserReq
	r, err := req.Bind(c)
	if err != nil {
		response.InvalidError(c, "Invalid request parameters")
		return
	}
	m := r.Model()

	user, err := h.userService.Update(m)
	if err != nil {
		response.HandleGormError(c, err)
		return
	}

	response.SuccessResponse[*model.User](c, "Updated", user)
}

// @Summary Delete a user
// @Description Deletes a user based on the provided user ID.
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID" example("67ae81a6-5284-436a-a2e5-54c3ebeaa241")
// @Success 200 {object} response.Response[any]
// @Failure 500 {object} response.Response[any]
// @Failure 400 {object} response.Response[any]
// @Failure 404 {object} response.Response[any]
// @Router /users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	var req request.GUserReq
	r, err := req.Bind(c)
	if err != nil {
		response.InvalidError(c, "Invalid request parameters")
		return
	}

	ok, err := h.userService.Delete(r)
	if err != nil || !ok {
		response.HandleGormError(c, err)
		return

	}

	response.SuccessResponse[interface{}](c, "Deleted", nil)
}
