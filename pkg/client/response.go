package client

import (
	"github.com/gin-gonic/gin"
)

// Convenience functions for direct Gin usage without HttpContext
func JsonResponse(c *gin.Context, status int32, message string, data any, err error) {
	response := HttpResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}

	if err != nil {
		response.Error = err.Error()
	} else {
		response.Error = ""
	}

	c.JSON(int(status), response)
}

func Success(c *gin.Context, message string, data any) {
	JsonResponse(c, 200, message, data, nil)
}

func Created(c *gin.Context, message string, data []any) {
	JsonResponse(c, 201, message, data, nil)
}

func Error(c *gin.Context, message string, err error) {
	JsonResponse(c, 400, message, nil, err)
}

func Unauthorized(c *gin.Context, err error) {
	JsonResponse(c, 401, "Unauthorized", nil, err)
}

func Forbidden(c *gin.Context, err error) {
	JsonResponse(c, 403, "Forbidden", nil, err)
}

func NotFound(c *gin.Context, err error) {
	JsonResponse(c, 404, "NotFound", nil, err)
}

func InternalServerError(c *gin.Context, err error) {
	JsonResponse(c, 500, "Internal Server Error", nil, err)
}
