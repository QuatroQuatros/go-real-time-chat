package chat

import (
	"log"

	"github.com/QuatroQuatros/go-real-time-chat/internal/domain"
)

type Hub struct {
	Connections map[*Connection]bool
	Broadcast   chan *domain.Message
	Register    chan *Connection
	Unregister  chan *Connection
}

func NewHub() *Hub {
	return &Hub{
		Connections: make(map[*Connection]bool),
		Broadcast:   make(chan *domain.Message),
		Register:    make(chan *Connection),
		Unregister:  make(chan *Connection),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case conn := <-h.Register:
			h.Connections[conn] = true
			log.Printf("🟢 User connected: %s", conn.User.Username)
		case conn := <-h.Unregister:
			if _, ok := h.Connections[conn]; ok {
				delete(h.Connections, conn)
				close(conn.Send)
				log.Printf("🔴 User disconnected: %s", conn.User.Username)
			}
		case msg := <-h.Broadcast:
			for conn := range h.Connections {
				select {
				case conn.Send <- msg:
				default:
					close(conn.Send)
					delete(h.Connections, conn)
				}
			}
		}
	}
}
