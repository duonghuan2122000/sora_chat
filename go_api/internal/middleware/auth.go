package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sora_chat/internal/consts"
	"sora_chat/internal/database"
	"sora_chat/internal/response"
	"sora_chat/internal/service"
	"strings"

	"github.com/gin-gonic/gin"
)

func getTokenInfoFromRedis(key string) (*service.TokenInfo, error) {
	data := database.RedisDB.Get(context.Background(), "sora:"+key)
	if err := data.Err(); err != nil {
		return nil, err
	}
	fmt.Println(data.Val())
	var tokenInfo *service.TokenInfo
	if err := json.Unmarshal([]byte(data.Val()), &tokenInfo); err != nil {
		return nil, err
	}
	return tokenInfo, nil
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			response.Error(c, http.StatusUnauthorized, consts.ErrUnauth, "Chưa xác thực", nil)
			return
		}

		token := parts[1]
		c.Set("token", token)
		tokenInfo, err := getTokenInfoFromRedis(token)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, consts.ErrUnauth, "Chưa xác thực", nil)
			return
		}
		c.Set("tokenInfo", tokenInfo)
		c.Next()
	}
}
