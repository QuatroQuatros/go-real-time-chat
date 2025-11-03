package chat

import (
	"log"
	"net/http"

	"github.com/QuatroQuatros/go-real-time-chat/infra/repository"
	"github.com/QuatroQuatros/go-real-time-chat/internal/domain"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // para testes locais; depois restrinja
	},
}

func ServeWs(hub *Hub, user *domain.User, msgRepo repository.MessageRepository, w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("❌ Failed to upgrade:", err)
		return
	}

	conn := &Connection{
		Ws:   ws,
		Send: make(chan *domain.Message),
		User: user,
	}

	hub.Register <- conn

	// Leitura de mensagens
	go func() {
		defer func() {
			hub.Unregister <- conn
			ws.Close()
		}()
		for {
			var msg domain.Message
			if err := ws.ReadJSON(&msg); err != nil {
				log.Println("❌ Read error:", err)
				break
			}

			// Adiciona usuário e timestamp
			msg.UserID = user.ID
			msg.CreatedAt = msg.CreatedAt.UTC()

			// Salva no banco
			if err := msgRepo.Create(&msg); err != nil {
				log.Println("❌ Failed to save message:", err)
				continue
			}

			hub.Broadcast <- &msg
		}
	}()

	// Escrita de mensagens
	go func() {
		for message := range conn.Send {
			if err := ws.WriteJSON(message); err != nil {
				log.Println("❌ Write error:", err)
				break
			}
		}
	}()
}
