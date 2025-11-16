package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/QuatroQuatros/go-real-time-chat/config"
	"github.com/QuatroQuatros/go-real-time-chat/infra/db"
	"github.com/QuatroQuatros/go-real-time-chat/internal/chat"
	"github.com/QuatroQuatros/go-real-time-chat/internal/domain"
	"github.com/QuatroQuatros/go-real-time-chat/internal/repository"
	"github.com/QuatroQuatros/go-real-time-chat/web"
)

func main() {
	if err := config.LoadEnv(); err != nil {
		log.Fatal(err)
	}

	db.Connect()

	msgRepo := repository.NewMessageRepository(db.DB)
	roomRepo := repository.NewRoomRepository(db.DB)

	hub := chat.NewHub()

	mux := http.NewServeMux()
	// ------------------------------
	// 🔥 Serve arquivos estáticos
	// ------------------------------
	mux.Handle("/", http.FileServer(http.FS(web.StaticFS)))

	// ------------------------------
	// Rotas da API
	// ------------------------------

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		username := r.URL.Query().Get("username")
		if username == "" {
			username = "Guest"
		}

		roomParam := r.URL.Query().Get("room")
		if roomParam == "" {
			roomParam = "general"
		}

		roomMap := map[string]uint{
			"general": 1,
			"random":  2,
			"games":   3,
			"support": 4,
		}

		roomID, ok := roomMap[roomParam]
		if !ok {
			roomID = 1 // fallback para "general"
		}

		// Cria usuário temporário (sem persistência) TODO: autenticação real
		user := &domain.User{
			ID:       1,
			Username: username,
		}

		chat.ServeWs(hub, user, roomID, msgRepo, w, r)
	})

	mux.HandleFunc("/rooms/", func(w http.ResponseWriter, r *http.Request) {
		// URL esperada: /rooms/{roomName}/messages
		path := r.URL.Path
		parts := strings.Split(path, "/")
		if len(parts) < 4 || parts[3] != "messages" {
			http.NotFound(w, r)
			return
		}

		roomIDStr := parts[2]
		roomID, err := strconv.ParseUint(roomIDStr, 10, 32)
		if err != nil {
			http.Error(w, "ID da sala inválido", http.StatusBadRequest)
			return
		}

		room, err := roomRepo.GetByID(uint(roomID))
		if err != nil {
			http.Error(w, "Sala não encontrada", http.StatusNotFound)
			return
		}

		limitStr := r.URL.Query().Get("limit")
		limit := 50
		if limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err == nil {
				limit = l
			}
		}

		msgs, err := msgRepo.GetByRoom(room.ID)
		if err != nil {
			http.Error(w, "Erro ao buscar mensagens", http.StatusInternalServerError)
			return
		}

		if len(msgs) > limit {
			msgs = msgs[len(msgs)-limit:]
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(msgs)
	})

	addr := fmt.Sprintf(":%s", config.Env.ServerPort)
	log.Printf("🚀 Server running on %s", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("❌ Server failed: %v", err)
	}
}
