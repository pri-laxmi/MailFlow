package models

import (
	"time"
)

type Template struct {
	Id        uint      `gorm:"primaryKey" json:"id"`
	UserId    uint      `gorm:"not null;index" json:"user_id"`
	Name      string    `gorm:"not null" json:"name"`
	Subject   string    `gorm:"not null" json:"subject"`
	Body      string    `gorm:"not null" json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      User      `gorm:"foreignKey:UserID" json:"-"`
}
