package service

import (
	"context"
	"fmt"

	"med-vito/api-go/internal/domain"
	"med-vito/api-go/internal/repository"
)

// LogService — бизнес-слой над логами (пока тонкая прокладка, как Nest LogService).
type LogService struct {
	repo *repository.LogPG
}

func NewLogService(repo *repository.LogPG) *LogService {
	return &LogService{repo: repo}
}

func (s *LogService) FindAll(ctx context.Context) ([]domain.Log, error) {
	items, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("log findAll: %w", err)
	}
	return items, nil
}
