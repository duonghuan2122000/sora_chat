package handler

import (
	"fmt"
	"net/http"
	"sora_chat/internal/ws"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// cấu hình upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var hub = ws.NewHub()

func init() {
	go hub.Run() // chạy goroutine quản lý client
}

// Websocket handler
func WebSocketHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}

	client := &ws.Client{
		ID:   uuid.New().String(),
		Conn: conn,
		Send: make(chan []byte),
	}

	hub.Register <- client
	go client.Read(hub)
	go client.Write()
}
