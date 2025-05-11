package main

import (
	"sora_chat/config"
	"sora_chat/internal/database"
	"sora_chat/internal/router"
)

func main() {
	// Load cấu hình từ file
	config.LoadConfig(".")

	database.ConnectRedis(config.AppConfig.RedisConn, config.AppConfig.RedisPass)

	hostName := config.AppConfig.HostName
	port := config.AppConfig.AppPort
	r := router.SetupRouter(config.AppConfig.TrustProxies, config.AppConfig)
	r.Run(hostName + ":" + port)
}
