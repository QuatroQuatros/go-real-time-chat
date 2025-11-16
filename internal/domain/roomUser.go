package domain

import "time"

type RoomUser struct {
	RoomID   uint
	UserID   uint
	JoinedAt time.Time
}
