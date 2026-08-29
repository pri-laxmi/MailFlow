package models

import "time"

type Job struct {
	ID         uint `gorm:"primaryKey" json:"id"`
	CampaignID uint `gorm:"not null;index" json:"campaign_id"`
	ContactID  uint `gorm:"not null;index" json:"contact_id"`

	Status       string `gorm:"not null;default:'pending'" json:"status"`
	RetryCount   int    `gorm:"not null;default:0" json:"retry_count"`
	ErrorMessage string `gorm:"type:text" json:"error_message,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Campaign Campaign `gorm:"foreignKey:CampaignID" json:"-"`
	Contact  Contact  `gorm:"foreignKey:ContactID" json:"-"`
}
