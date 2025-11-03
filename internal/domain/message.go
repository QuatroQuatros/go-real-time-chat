package domain

import "time"

type Message struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"not null"`
	RoomID    uint   `gorm:"not null"`
	Content   string `gorm:"type:text;not null"`
	CreatedAt time.Time

	User User `gorm:"foreignKey:UserID"`
	Room Room `gorm:"foreignKey:RoomID"`
}
