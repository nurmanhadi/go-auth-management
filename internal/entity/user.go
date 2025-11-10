package entity

import (
	"auth-management/pkg/enum"
	"time"
)

type User struct {
	Id        string    `gorm:"primaryKey"`
	Username  string    `gorm:"unique; not null"`
	Password  string    `gorm:"not null"`
	Role      enum.ROLE `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
