// router.go

package router

import (
	"sora_chat/config"
	"sora_chat/internal/database"
	"sora_chat/internal/handler"
	"sora_chat/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(trustProxies []string, appConfig config.Config) *gin.Engine {
	r := gin.Default()

	r.SetTrustedProxies(trustProxies)

	mongoDb := database.ConnectMongo(appConfig.MongoURI, appConfig.MongoDB)
	handler.InitHandlerUser(mongoDb)

	api := r.Group("/api")
	{
		api.GET("/healthz", handler.Healthz)

		// API cấp quyền truy cập
		api.POST("/connect/token", handler.GrantToken)

		api.GET("/users", handler.GetUsers)
		api.GET("/user/cu", middleware.AuthMiddleware(), handler.GetCurrentUserInfo)
		api.POST("/users", handler.CreateUser)

		api.GET("/ws", handler.WebSocketHandler)
	}

	return r
}
