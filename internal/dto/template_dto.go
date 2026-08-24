package dto

import (
	"time"
)

type CreateTemplateRequest struct {
	Name    string `json:"name" binding:"required"`
	Subject string `json:"subject" binding:"required"`
	Body    string `json:"body" binding:"required"`
}
type UpdateTemplateRequest struct {
	Name    *string `json:"name"`
	Subject *string `json:"subject"`
	Body    *string `json:"body"`
}

type TemplateResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Subject   string    `json:"subject"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
