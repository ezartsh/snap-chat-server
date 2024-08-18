package models

import (
	"time"
)

type Room struct {
	ID        *uint
	RoomUID   string
	RoomType  string
	Name      string
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

func (m Room) TableName() string {
	return "rooms"
}
