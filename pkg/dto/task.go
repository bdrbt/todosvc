package dto

import (
	"time"

	"github.com/google/uuid"
)

type TaskDTO struct {
	ID         uuid.UUID `json:"id"`
	Decription string    `json:"description"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
