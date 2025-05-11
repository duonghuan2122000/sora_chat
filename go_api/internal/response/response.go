package response

import "github.com/gin-gonic/gin"

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
}

type APIError struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// trả về thành công
func Success(c *gin.Context, status int, data interface{}, message string) {
	c.JSON(status, APIResponse{
		Success: true,
		Data:    data,
		Message: message,
	})
}

// trả về lỗi
func Error(c *gin.Context, status int, code string, message string, details interface{}) {
	c.AbortWithStatusJSON(status, APIResponse{
		Success: false,
		Error: &APIError{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}
