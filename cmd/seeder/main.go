package main

import (
	"log"
	"time"

	"github.com/QuatroQuatros/go-real-time-chat/config"
	"github.com/QuatroQuatros/go-real-time-chat/infra/db"
	"github.com/QuatroQuatros/go-real-time-chat/internal/domain"
	"github.com/QuatroQuatros/go-real-time-chat/internal/repository"
	"github.com/QuatroQuatros/go-real-time-chat/internal/shared"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Carrega variáveis de ambiente
	if err := config.LoadEnv(); err != nil {
		log.Fatal(err)
	}

	// Conecta ao banco
	db.Connect()

	// Cria repositórios
	userRepo := repository.NewUserRepository(db.DB)
	roomRepo := repository.NewRoomRepository(db.DB)

	// Seed Users
	users := []domain.User{
		{
			Username:  "TestUser",
			Password:  hash("123456"),
			Guest:     false,
			LastLogin: time.Now(),
			AuditInfo: shared.AuditInfo{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
		{
			Username:  "Alice",
			Password:  hash("senhaAlice"),
			Guest:     false,
			LastLogin: time.Now(),
			AuditInfo: shared.AuditInfo{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
		{
			Username:  "Bob",
			Password:  hash("senhaBob"),
			Guest:     false,
			LastLogin: time.Now(),
			AuditInfo: shared.AuditInfo{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
	}

	for _, u := range users {
		if err := userRepo.Create(&u); err != nil {
			log.Printf("⚠️ Usuário já existe ou erro: %v", err)
		} else {
			log.Printf("✅ Usuário criado: %s", u.Username)
		}
	}

	// Seed Rooms
	rooms := []domain.Room{
		{Name: "Geral", CreatedAt: time.Now()},
		{Name: "Random", CreatedAt: time.Now()},
		{Name: "Games", CreatedAt: time.Now()},
		{Name: "Support", CreatedAt: time.Now()},
	}

	for _, r := range rooms {
		if err := roomRepo.Create(&r); err != nil {
			log.Printf("⚠️ Sala já existe ou erro: %v", err)
		} else {
			log.Printf("✅ Sala criada: %s", r.Name)
		}
	}

	log.Println("🎉 Seeder finalizado!")
}

func hash(pwd string) string {
	h, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Erro ao fazer hash: %v", err)
	}
	return string(h)
}
