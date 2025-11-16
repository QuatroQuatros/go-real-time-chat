package domain

import "time"

type Room struct {
	ID        uint    `gorm:"primaryKey"`
	Name      string  `gorm:"unique;not null"`
	IsPrivate bool    `gorm:"default:false"` // true = chat privado
	OwnerID   uint    // usuário que criou a sala (opcional)
	Users     []*User `gorm:"many2many:room_users;"` // relação N:N
	CreatedAt time.Time
}
