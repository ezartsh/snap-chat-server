package models

import (
	"time"
)

type RoomAudience struct {
	ID        *uint
	RoomID    uint
	UserID    uint
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

func (m RoomAudience) TableName() string {
	return "rooms_audiences"
}
