package dto

import (
	"time"
)

type CreateContactRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}
type UpdateContactRequest struct {
	Name  *string `json:"name" binding:"required"`
	Email *string `json:"email" binding:"required,email"`
}
type ContactResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
