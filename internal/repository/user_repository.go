package repository

import (
	"github.com/QuatroQuatros/go-real-time-chat/infra/repository"
	"github.com/QuatroQuatros/go-real-time-chat/internal/domain"
	"gorm.io/gorm"
)

type userRepositoryGorm struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepositoryGorm{db: db}
}

func (r *userRepositoryGorm) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *userRepositoryGorm) GetByID(id uint) (*domain.User, error) {
	var user domain.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryGorm) GetByUsername(username string) (*domain.User, error) {
	var user domain.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
