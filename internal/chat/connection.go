package chat

import (
	"github.com/QuatroQuatros/go-real-time-chat/internal/domain"
	"github.com/gorilla/websocket"
)

type Connection struct {
	// websocket connection
	Ws   *websocket.Conn
	Send chan *domain.Message
	User *domain.User
}
