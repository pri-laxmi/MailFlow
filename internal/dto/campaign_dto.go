package dto

import "time"

type CreateCampaignRequest struct {
	Name        string     `json:"name" binding:"required"`
	TemplateID  uint       `json:"template_id" binding:"required"`
	ScheduledAt *time.Time `json:"scheduled_at"`
}

type UpdateCampaignRequest struct {
	Name        *string    `json:"name"`
	TemplateID  *uint      `json:"template_id"`
	Status      *string    `json:"status"`
	ScheduledAt *time.Time `json:"scheduled_at"`
}
