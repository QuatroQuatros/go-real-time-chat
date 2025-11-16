package repository

import (
	"github.com/QuatroQuatros/go-real-time-chat/infra/repository"
	"github.com/QuatroQuatros/go-real-time-chat/internal/domain"
	"github.com/google/uuid"
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

func (r *roomRepositoryGorm) GetByUser(userID uint) ([]*domain.Room, error) {
	var rooms []*domain.Room
	if err := r.db.Joins("JOIN room_users ru ON ru.room_id = rooms.id").
		Where("ru.user_id = ?", userID).
		Preload("Users").
		Find(&rooms).Error; err != nil {
		return nil, err
	}
	return rooms, nil
}

func (r *roomRepositoryGorm) CreatePrivateRoom(user1ID, user2ID uint) (*domain.Room, error) {
	room := &domain.Room{
		Name:      uuid.New().String(),
		IsPrivate: true,
		Users: []*domain.User{
			{ID: user1ID},
			{ID: user2ID},
		},
	}

	if err := r.db.Create(room).Error; err != nil {
		return nil, err
	}
	return room, nil
}

func (r *roomRepositoryGorm) AddUserToRoom(roomID, userID uint) error {
	var room domain.Room
	if err := r.db.Preload("Users").First(&room, roomID).Error; err != nil {
		return err
	}

	var user domain.User
	if err := r.db.First(&user, userID).Error; err != nil {
		return err
	}

	return r.db.Model(&room).Association("Users").Append(&user)
}

func (r *roomRepositoryGorm) RemoveUserFromRoom(roomID, userID uint) error {
	var room domain.Room
	if err := r.db.Preload("Users").First(&room, roomID).Error; err != nil {
		return err
	}

	var user domain.User
	if err := r.db.First(&user, userID).Error; err != nil {
		return err
	}

	return r.db.Model(&room).Association("Users").Delete(&user)
}
