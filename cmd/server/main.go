package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/QuatroQuatros/go-real-time-chat/config"
	"github.com/QuatroQuatros/go-real-time-chat/infra/db"
	"github.com/QuatroQuatros/go-real-time-chat/internal/chat"
	"github.com/QuatroQuatros/go-real-time-chat/internal/domain"
	"github.com/QuatroQuatros/go-real-time-chat/internal/repository"
	"github.com/QuatroQuatros/go-real-time-chat/web"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func main() {
	if err := config.LoadEnv(); err != nil {
		log.Fatal(err)
	}

	db.Connect()

	msgRepo := repository.NewMessageRepository(db.DB)
	roomRepo := repository.NewRoomRepository(db.DB)

	hub := chat.NewHub()

	r := gin.Default()

	// ------------------------------
	// 🔥 Serve arquivos estáticos
	// ------------------------------

	fs, err := static.EmbedFolder(web.StaticFS, ".")
	if err != nil {
		log.Fatal(err)
	}

	r.Use(static.Serve("/", fs))

	r.NoRoute(func(c *gin.Context) {
		fmt.Printf("%s doesn't exists, redirect on /\n", c.Request.URL.Path)
		c.Redirect(http.StatusMovedPermanently, "/")
	})

	// ------------------------------
	// Rotas da API
	// ------------------------------

	api := r.Group("/api")

	api.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	api.GET("/ws", func(c *gin.Context) {
		username := c.Query("username")
		if username == "" {
			username = "Guest"
		}

		roomParam := c.Query("room")
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

		chat.ServeWs(hub, user, roomID, msgRepo, c.Writer, c.Request)
	})

	api.GET("/rooms/:roomID/messages", func(c *gin.Context) {
		roomIDStr := c.Param("roomID")

		roomID, err := strconv.ParseUint(roomIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"eror": "ID da sala inválido"})
			return
		}

		room, err := roomRepo.GetByID(uint(roomID))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"eror": "Sala não encontrada"})
			return
		}

		limitStr := c.Query("limit")
		limit := 50
		if limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err == nil {
				limit = l
			}
		}

		msgs, err := msgRepo.GetByRoom(room.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"eror": "Erro ao buscar mensagens"})
			return
		}

		if len(msgs) > limit {
			msgs = msgs[len(msgs)-limit:]
		}

		c.JSON(http.StatusOK, msgs)
	})

	addr := fmt.Sprintf(":%s", config.Env.ServerPort)
	log.Printf("🚀 Server running on %s", addr)

	if err := r.Run(addr); err != nil {
		log.Fatalf("❌ Server failed: %v", err)
	}
}
