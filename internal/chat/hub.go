package chat

import (
	"log"

	"github.com/QuatroQuatros/go-real-time-chat/internal/domain"
)

type RoomHub struct {
	Connections map[*Connection]bool
	Broadcast   chan *domain.Message
	Register    chan *Connection
	Unregister  chan *Connection
}

func NewRoomHub() *RoomHub {
	return &RoomHub{
		Connections: make(map[*Connection]bool),
		Broadcast:   make(chan *domain.Message),
		Register:    make(chan *Connection),
		Unregister:  make(chan *Connection),
	}
}

type Hub struct {
	Rooms map[uint]*RoomHub // roomID -> RoomHub
}

func NewHub() *Hub {
	return &Hub{
		Rooms: make(map[uint]*RoomHub),
	}
}

func (h *Hub) GetRoomHub(roomID uint) *RoomHub {
	if hub, ok := h.Rooms[roomID]; ok {
		return hub
	}

	hub := NewRoomHub()
	h.Rooms[roomID] = hub
	go hub.Run()
	return hub
}

func (rh *RoomHub) Run() {
	for {
		select {
		case conn := <-rh.Register:
			rh.Connections[conn] = true
			log.Printf("🟢 User connected: %s", conn.User.Username)
		case conn := <-rh.Unregister:
			if _, ok := rh.Connections[conn]; ok {
				delete(rh.Connections, conn)
				close(conn.Send)
				log.Printf("🔴 User disconnected: %s", conn.User.Username)
			}
		case msg := <-rh.Broadcast:
			for conn := range rh.Connections {
				select {
				case conn.Send <- msg:
				default:
					close(conn.Send)
					delete(rh.Connections, conn)
				}
			}
		}
	}
}
