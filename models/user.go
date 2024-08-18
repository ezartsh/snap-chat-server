package models

import "time"

type User struct {
	ID        *uint
	Name      string
	Username  string
	Password  string
	CreatedAt *time.Time
	UpdatedAt *time.Time
}
