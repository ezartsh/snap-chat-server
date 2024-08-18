package models

import (
	"time"
)

type Group struct {
	ID        *uint
	RoomID    int
	Key       string
	Name      string
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

func (m Group) TableName() string {
	return "groups"
}
