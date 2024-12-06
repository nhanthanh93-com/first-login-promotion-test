package response

import (
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
	"math"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Pagination struct {
	CurrentPage int32 `json:"current_page"`
	PerPage     int32 `json:"per_page"`
	TotalItem   int64 `json:"total_item"`
	TotalPage   int32 `json:"total_page"`
}

type Response[T any] struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    *T     `json:"data,omitempty"`
}

type ResponseData[T any] struct {
	Pagination Pagination `json:"pagination"`
	Results    *T         `json:"results"`
}

func PaginationResponse(page, limit, perPage int32, totalItems int64) Pagination {
	totalPage := int32(math.Ceil(float64(totalItems) / float64(limit)))

	return Pagination{
		CurrentPage: page,
		PerPage:     perPage,
		TotalPage:   totalPage,
		TotalItem:   totalItems,
	}
}

func mapStatusText(statusCode int) string {
	return strings.ReplaceAll(strings.ToUpper(http.StatusText(statusCode)), " ", "_")
}

func SuccessResponse[T any](c *gin.Context, message string, data T) {
	c.JSON(http.StatusOK, Response[T]{
		Status:  "OK",
		Message: message,
		Data:    &data,
	})
}

func ErrorResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response[interface{}]{
		Status:  mapStatusText(statusCode),
		Message: message,
		Data:    &data,
	})
}

func HandleGormError(c *gin.Context, err error) {
	if err == nil {
		return
	}
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		NotFoundError(c, "Record not found")
		return
	case errors.Is(err, gorm.ErrInvalidData):
		InvalidError(c, "Invalid data: "+err.Error())
		return
	case errors.Is(err, gorm.ErrDuplicatedKey):
		ExistedError(c, "Duplicated key: "+err.Error())
		return
	case errors.Is(err, gorm.ErrPrimaryKeyRequired):
		InvalidError(c, "Primary key required: "+err.Error())
		return
	case errors.Is(err, gorm.ErrEmptySlice):
		InvalidError(c, "Empty input slice")
		return
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			ExistedError(c, "Duplicated key: "+err.Error())
			return
		}
	}

	GeneralError(c, "Database error: "+err.Error())
}

func UnauthorizedError(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusUnauthorized, message, nil)
}

func InvalidError(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusBadRequest, message, nil)
}

func ForbiddenError(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusForbidden, message, nil)
}

func ExistedError(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusConflict, message, nil)
}

func NotFoundError(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusNotFound, message, nil)
}

func GeneralError(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusInternalServerError, message, nil)
}

func BindError(c *gin.Context, statusCode int, message string) {
	ErrorResponse(c, statusCode, message, nil)
}
