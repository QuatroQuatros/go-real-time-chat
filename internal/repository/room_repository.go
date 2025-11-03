package repository

import (
	"github.com/QuatroQuatros/go-real-time-chat/infra/repository"
	"github.com/QuatroQuatros/go-real-time-chat/internal/domain"
	"gorm.io/gorm"
)

// Implementação concreta
type roomRepositoryGorm struct {
	db *gorm.DB
}

// Construtor
func NewRoomRepository(db *gorm.DB) repository.RoomRepository {
	return &roomRepositoryGorm{db: db}
}

func (r *roomRepositoryGorm) Create(room *domain.Room) error {
	return r.db.Create(room).Error
}

func (r *roomRepositoryGorm) GetByID(id uint) (*domain.Room, error) {
	var room domain.Room
	if err := r.db.First(&room, id).Error; err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *roomRepositoryGorm) GetByName(name string) (*domain.Room, error) {
	var room domain.Room
	if err := r.db.Where("name = ?", name).First(&room).Error; err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *roomRepositoryGorm) GetAll() ([]*domain.Room, error) {
	var rooms []*domain.Room
	if err := r.db.Find(&rooms).Error; err != nil {
		return nil, err
	}
	return rooms, nil
}
