package ws

import "github.com/gorilla/websocket"

func (c *Client) Read(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		// khi nhận được message -> broadcast tới tất cả
		hub.Broadcast <- message
	}
}

func (c *Client) Write() {
	for msg := range c.Send {
		c.Conn.WriteMessage(websocket.TextMessage, msg)
	}
}
