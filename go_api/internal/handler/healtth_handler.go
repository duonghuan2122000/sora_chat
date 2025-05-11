package handler

import (
	"net/http"
	"sora_chat/internal/response"

	"github.com/gin-gonic/gin"
)

// hàm trả về kết quả healthz cho app
func Healthz(c *gin.Context) {
	response.Success(c, http.StatusOK, nil, "App is running")
}
