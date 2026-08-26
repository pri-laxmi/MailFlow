package models

import "time"

type Schedule struct {
	draft string `json:"draft"`
}
type Campaign struct {
	ID         uint `gorm:"primaryKey" json:"id"`
	UserID     uint `gorm:"not null;index" json:"user_id"`
	TemplateID uint `gorm:"not null;index" json:"template_id"`

	Name        string     `gorm:"not null" json:"name"`
	Status      string     `gorm:"not null;default:'draft'" json:"status"`
	ScheduledAt *time.Time `json:"scheduled_at"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	User     User     `gorm:"foreignKey:UserID" json:"-"`
	Template Template `gorm:"foreignKey:TemplateID" json:"template,omitempty"`
}
