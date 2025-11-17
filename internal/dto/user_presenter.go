package dto

import (
	"time"

	"github.com/QuatroQuatros/go-real-time-chat/internal/domain"
)

type UserPresenter struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Guest     bool      `json:"guest"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AuthPresenter struct {
	User  *UserPresenter `json:"user"`
	Token string         `json:"token"`
}

func NewUserPresenter(user *domain.User) *UserPresenter {
	if user == nil {
		return nil
	}

	return &UserPresenter{
		ID:        user.ID,
		Username:  user.Username,
		Guest:     user.Guest,
		CreatedAt: user.AuditInfo.CreatedAt,
		UpdatedAt: user.AuditInfo.UpdatedAt,
	}
}
