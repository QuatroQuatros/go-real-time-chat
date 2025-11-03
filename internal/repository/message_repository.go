package repository

import (
	"github.com/QuatroQuatros/go-real-time-chat/infra/repository"
	"github.com/QuatroQuatros/go-real-time-chat/internal/domain"
	"gorm.io/gorm"
)

// Implementação concreta
type messageRepositoryGorm struct {
	db *gorm.DB
}

// Construtor
func NewMessageRepository(db *gorm.DB) repository.MessageRepository {
	return &messageRepositoryGorm{db: db}
}

func (r *messageRepositoryGorm) Create(message *domain.Message) error {
	return r.db.Create(message).Error
}

func (r *messageRepositoryGorm) GetByID(id uint) (*domain.Message, error) {
	var msg domain.Message
	if err := r.db.Preload("User").Preload("Room").First(&msg, id).Error; err != nil {
		return nil, err
	}
	return &msg, nil
}

func (r *messageRepositoryGorm) GetByRoom(roomID uint) ([]*domain.Message, error) {
	var msgs []*domain.Message
	if err := r.db.Preload("User").Where("room_id = ?", roomID).Order("created_at asc").Find(&msgs).Error; err != nil {
		return nil, err
	}
	return msgs, nil
}
