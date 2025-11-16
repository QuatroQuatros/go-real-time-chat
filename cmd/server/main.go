package main

import (
	"fmt"
	"log"
	"net/http"

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

	hub := chat.NewHub()
	go hub.Run()

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
		// Para testes, cria usuário fake
		user := &domain.User{ID: 1, Username: "TestUser"}
		chat.ServeWs(hub, user, msgRepo, w, r)
	})

	addr := fmt.Sprintf(":%s", config.Env.ServerPort)
	log.Printf("🚀 Server running on %s", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("❌ Server failed: %v", err)
	}
}
