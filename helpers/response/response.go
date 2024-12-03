package response

import (
	"errors"
	"gorm.io/gorm"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func mapStatusText(statusCode int) string {
	return strings.ReplaceAll(strings.ToUpper(http.StatusText(statusCode)), " ", "_")
}

func SuccessResponse(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Status:  "OK",
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Status:  mapStatusText(statusCode),
		Message: message,
		Data:    data,
	})
}

func HandleGormError(c *gin.Context, err error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		NotFoundError(c, "Record not found")
	} else if errors.Is(err, gorm.ErrInvalidData) {
		InvalidError(c, "Invalid data: "+err.Error())
	} else if errors.Is(err, gorm.ErrDuplicatedKey) {
		ExistedError(c, "Duplicated key: "+err.Error())
	} else if errors.Is(err, gorm.ErrPrimaryKeyRequired) {
		InvalidError(c, "Primary key required: "+err.Error())
	} else if errors.Is(err, gorm.ErrEmptySlice) {
		InvalidError(c, "Empty input slice")
	} else {
		GeneralError(c, "Database error: "+err.Error())
	}
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
