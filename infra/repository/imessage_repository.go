package repository

import "github.com/QuatroQuatros/go-real-time-chat/internal/domain"

// Interface define as operações para mensagens
type MessageRepository interface {
	Create(message *domain.Message) error
	GetByID(id uint) (*domain.Message, error)
	GetByRoom(roomID uint) ([]*domain.Message, error)
}
