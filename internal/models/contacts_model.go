package models

import (
	"time"
)

type Contact struct {
	id        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	userId    uint      `gorm:"not null" json:"user_id"`
	Name      string    `gorm:"not null" json:"name"`
	Email     string    `gorm:"unique;not null" json:"email"`
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`
}
