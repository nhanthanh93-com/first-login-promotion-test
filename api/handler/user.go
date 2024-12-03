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

// CreateTodo godoc
// @Summary Lấy thông tin người dùng
// @Description API lấy thông tin người dùng qua ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} User
// @Failure 404 {object} ErrorResponse
// @Router /users/{id} [get]
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

	response.SuccessResponse(c, "Created", user)
}

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

	response.SuccessResponse(c, "Created", user)
}

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

	response.SuccessResponse(c, "Success", user)
}

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

	response.SuccessResponse(c, "List", gin.H{
		"results":    rs,
		"pagination": r.FormatPagination(count, int64(len(rs))),
	})

}

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

	response.SuccessResponse(c, "Updated", user)
}

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

	response.SuccessResponse(c, "Deleted", nil)
}
