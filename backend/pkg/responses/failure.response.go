package responses

import "github.com/gin-gonic/gin"

type ApiFailureResponse struct {
	Success      bool        `json:"success"`
	ErrorMessage string      `json:"error_message"`
	Error        interface{} `json:"error,omitempty"`
}

func ApiFailure(ctx *gin.Context, errorMessage string, statusCode int, err interface{}) {
	response := &ApiFailureResponse{
		Success:      false,
		ErrorMessage: errorMessage,
		Error:        err,
	}

	ctx.JSON(statusCode, response)
}
