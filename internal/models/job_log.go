package models

import "time"

type JobLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	JobID     uint      `gorm:"index" json:"job_id"`
	Event     string    `gorm:"not null" json:"event"`
	Message   string    `gorm:"type:text" json:"message"`
	CreatedAt time.Time `json:"created_at"`

	Job Job `gorm:"foreignKey:JobID" json:"job"`
}
