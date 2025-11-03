package repository

import "github.com/QuatroQuatros/go-real-time-chat/internal/domain"

type UserRepository interface {
	Create(user *domain.User) error
	GetByID(id uint) (*domain.User, error)
	GetByUsername(username string) (*domain.User, error)
}
