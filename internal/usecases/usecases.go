package usecases

import (
	"github.com/bdrbt/todo/internal/repository"
	"go.uber.org/zap"
)

type UC struct {
	logger *zap.Logger
	repo   *repository.Repository
}

func New(r *repository.Repository, l *zap.Logger) *UC {
	uc := &UC{
		logger: l,
		repo:   r,
	}

	return uc
}
