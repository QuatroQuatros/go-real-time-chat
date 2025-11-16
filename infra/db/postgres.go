package db

import (
	"fmt"
	"log"
	"os"

	"github.com/QuatroQuatros/go-real-time-chat/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	host := os.Getenv("DB_SERVER")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=America/Sao_Paulo",
		host, user, password, dbname, port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	DB = db
	log.Println("✅ Connected to database")

	// Auto-migration
	err = db.AutoMigrate(&domain.User{}, &domain.Room{}, &domain.Message{})
	if err != nil {
		log.Fatalf("❌ Failed to migrate tables: %v", err)
	}

	log.Println("✅ Database migrated successfully")
}
