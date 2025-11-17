package domain

import (
	"time"

	"github.com/QuatroQuatros/go-real-time-chat/internal/shared"
	userErrors "github.com/QuatroQuatros/go-real-time-chat/internal/shared/errors"
)

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"unique;not null"`
	Password  string `json:"-"`
	Guest     bool   `gorm:"default:false"`
	LastLogin time.Time
	AuditInfo shared.AuditInfo `gorm:"embedded"`
}

func NewUserFomInput(username, password string, guest bool) (*User, error) {
	u := &User{
		Username:  username,
		Password:  password,
		Guest:     guest,
		AuditInfo: shared.AuditInfo{CreatedAt: time.Now()},
	}

	if err := u.ValidateInternalState(); err != nil {
		return nil, err
	}

	return u, nil
}

func (u *User) ValidateInternalState() error {
	if u.hasUsername() && len(u.Username) < 3 {
		return userErrors.ErrUsernameLength
	}
	if u.hasPassword() && len(u.Password) < 8 {
		return userErrors.ErrPasswordLength
	}
	return nil
}

func IsGuest(guest bool) bool {
	return guest
}

func (u *User) hasUsername() bool {
	return u.Username != ""
}

func (u *User) hasPassword() bool {
	return u.Password != ""
}
