// Package handler

package handler

import (
	"errors"
	"net/http"
	"sora_chat/internal/consts"
	sora_errors "sora_chat/internal/errors"
	"sora_chat/internal/repository"
	"sora_chat/internal/response"
	"sora_chat/internal/service"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var userService service.UserService

func InitHandlerUser(db *mongo.Database) {
	userRepo := repository.NewUserRepository(db)
	userService = service.NewUserService(userRepo)
}

func GetUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "List of users"})
}

func CreateUser(c *gin.Context) {
	var payload service.CreateUserReqDto

	// Gắn JSON body vào struct
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, consts.ErrBadRequest, "Request Body không hợp lệ", nil)
		return
	}

	result, err := userService.CreateUser(payload)
	if err != nil {
		response.Error(c, http.StatusBadRequest, consts.ErrBadRequest, "Có lỗi khi cấp quyền", nil)
		return
	}
	response.Success(c, http.StatusOK, *result, "Thành công")
}

// Thực hiện cấp quyền truy cập
func GrantToken(c *gin.Context) {
	var payload service.LoginUserReqDto

	// Gắn JSON body vào struct
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, consts.ErrBadRequest, "Request Body không hợp lệ", nil)
		return
	}
	result, err := userService.GrantToken(payload)
	if err != nil {
		var logicErr *sora_errors.LogicError
		if errors.As(err, &logicErr) {
			response.Error(c, http.StatusBadRequest, logicErr.Code, logicErr.Message, nil)
			return
		}

		response.Error(c, http.StatusBadRequest, consts.ErrBadRequest, "Có lỗi khi cấp quyền", nil)
		return
	}
	response.Success(c, http.StatusOK, *result, "Thành công")
}

// lấy thông tin user của phiên hiện tại
func GetCurrentUserInfo(c *gin.Context) {
	v, _ := c.Get("token")
	token := v.(string)
	v, _ = c.Get("tokenInfo")
	tokenInfo := v.(*service.TokenInfo)
	result, err := userService.GetCurrentUserInfo(token, tokenInfo)
	if err != nil {
		response.Error(c, http.StatusBadRequest, consts.ErrBadRequest, "Có lỗi khi cấp quyền", nil)
		return
	}
	response.Success(c, http.StatusOK, *result, "Thành công")
}
