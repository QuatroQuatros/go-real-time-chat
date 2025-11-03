package repository

import "github.com/QuatroQuatros/go-real-time-chat/internal/domain"

// Interface define operações para salas
type RoomRepository interface {
	Create(room *domain.Room) error
	GetByID(id uint) (*domain.Room, error)
	GetByName(name string) (*domain.Room, error)
	GetAll() ([]*domain.Room, error)
}
