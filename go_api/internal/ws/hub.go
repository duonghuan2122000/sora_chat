package ws

import (
	"context"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
)

var ctx = context.Background()

type Client struct {
	ID   string
	Conn *websocket.Conn
	Send chan []byte
}

type Hub struct {
	Clients    map[string]*Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan []byte
	Adapter    *RedisAdapter
	mu         sync.Mutex
}

type RedisAdapter struct {
	Client  *redis.Client
	Channel string
}

func NewRedisAdapter() *RedisAdapter {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis host
	})
	return &RedisAdapter{
		Client:  rdb,
		Channel: "chat",
	}
}

func (r *RedisAdapter) Publish(message []byte) {
	r.Client.Publish(ctx, r.Channel, message)
}

func (r *RedisAdapter) Subscribe(hub *Hub) {
	sub := r.Client.Subscribe(ctx, r.Channel)
	ch := sub.Channel()
	go func() {
		for msg := range ch {
			hub.Broadcast <- []byte(msg.Payload)
		}
	}()
}

func NewHub() *Hub {
	adapter := NewRedisAdapter()
	h := &Hub{
		Clients:    make(map[string]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan []byte),
		Adapter:    adapter,
	}
	adapter.Subscribe(h)
	return h
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client.ID] = client
			h.mu.Unlock()

		case client := <-h.Unregister:
			h.mu.Lock()
			delete(h.Clients, client.ID)
			close(client.Send)
			h.mu.Unlock()

		case message := <-h.Broadcast:
			// Gửi message ra Redis (cho các instance khác)
			h.Adapter.Publish(message)

			// Gửi local client
			h.mu.Lock()
			for _, client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client.ID)
				}
			}
			h.mu.Unlock()
		}

	}
}
