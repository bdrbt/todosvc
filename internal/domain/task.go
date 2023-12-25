package domain

import (
	"time"

	"github.com/bdrbt/todo/pkg/dto"
	"github.com/google/uuid"
)

type Task struct {
	ID         uuid.UUID
	Decription string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  time.Time
}

func (t *Task) DTO() (*dto.TaskDTO, error) {
	dto := dto.TaskDTO{
		ID:         t.ID,
		Decription: t.Decription,
		CreatedAt:  t.CreatedAt,
		UpdatedAt:  t.UpdatedAt,
	}
	return &dto, nil
}
