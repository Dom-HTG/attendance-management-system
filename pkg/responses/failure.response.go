package responses

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ApiFailureResponse struct {
	Success      bool        `json:"success"`
	ErrorMessage string      `json:"error_message"`
	Error        interface{} `json:"error,omitempty"`
	Details      interface{} `json:"details,omitempty"`
}

func ApiFailure(ctx *gin.Context, errorMessage string, statusCode int, err interface{}) {
	response := &ApiFailureResponse{
		Success:      false,
		ErrorMessage: errorMessage,
	}

	// Handle validation errors specifically
	if err != nil {
		switch e := err.(type) {
		case validator.ValidationErrors:
			// Convert validation errors to a readable format
			validationErrors := make(map[string]string)
			for _, fieldErr := range e {
				validationErrors[fieldErr.Field()] = formatValidationError(fieldErr)
			}
			response.Details = validationErrors
			response.Error = "validation failed"
		case error:
			// Handle regular errors
			response.Error = e.Error()
		default:
			// For other types, convert to string
			response.Error = fmt.Sprintf("%v", err)
		}
	}

	ctx.JSON(statusCode, response)
}

func formatValidationError(fieldErr validator.FieldError) string {
	switch fieldErr.Tag() {
	case "required":
		return "this field is required"
	case "email":
		return "invalid email format"
	case "min":
		return fmt.Sprintf("minimum length is %s", fieldErr.Param())
	case "max":
		return fmt.Sprintf("maximum length is %s", fieldErr.Param())
	default:
		return fmt.Sprintf("validation failed on '%s'", fieldErr.Tag())
	}
}
