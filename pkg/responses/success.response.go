package responses

import "github.com/gin-gonic/gin"

type ApiSuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func ApiSuccess(ctx *gin.Context, statusCode int, message string, data interface{}) {
	response := &ApiSuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	}

	ctx.JSON(statusCode, response)

}
