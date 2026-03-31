package service

import (
	"context"
	"errors"
	"strings"

	"med-vito/api-go/internal/domain"
	"med-vito/api-go/internal/repository"
)

type KnowledgeBaseService struct {
	repo *repository.KnowledgeBasePG
}

func NewKnowledgeBaseService(repo *repository.KnowledgeBasePG) *KnowledgeBaseService {
	return &KnowledgeBaseService{repo: repo}
}

func (s *KnowledgeBaseService) CreateArticle(ctx context.Context, title, content string) (*domain.CreateKnowledgeBaseArticleResponse, error) {
	title = strings.TrimSpace(title)
	content = strings.TrimSpace(content)

	if title == "" {
		return nil, &AppError{Status: 400, Message: "Заголовок статьи обязателен для заполнения"}
	}
	if content == "" {
		return nil, &AppError{Status: 400, Message: "Текст статьи обязателен для заполнения"}
	}

	article, err := s.repo.CreateArticle(ctx, title, content)
	if err != nil {
		return nil, err
	}
	return &domain.CreateKnowledgeBaseArticleResponse{
		Message: "Статья успешно создана",
		Article: *article,
	}, nil
}

func (s *KnowledgeBaseService) ListArticles(ctx context.Context) ([]domain.KnowledgeBaseArticle, error) {
	return s.repo.ListArticles(ctx)
}

func (s *KnowledgeBaseService) GetArticleByID(ctx context.Context, id int32) (*domain.KnowledgeBaseArticle, error) {
	article, err := s.repo.GetArticleByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, &AppError{Status: 404, Message: "Статья не найдена"}
		}
		return nil, err
	}
	return article, nil
}

func (s *KnowledgeBaseService) UpdateArticle(ctx context.Context, id int32, title, content string) (*domain.UpdateKnowledgeBaseArticleResponse, error) {
	title = strings.TrimSpace(title)
	content = strings.TrimSpace(content)

	if title == "" {
		return nil, &AppError{Status: 400, Message: "Заголовок статьи обязателен для заполнения"}
	}
	if content == "" {
		return nil, &AppError{Status: 400, Message: "Текст статьи обязателен для заполнения"}
	}

	article, err := s.repo.UpdateArticle(ctx, id, title, content)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, &AppError{Status: 404, Message: "Статья не найдена"}
		}
		return nil, err
	}
	return &domain.UpdateKnowledgeBaseArticleResponse{
		Message: "Статья успешно обновлена",
		Article: *article,
	}, nil
}

func (s *KnowledgeBaseService) DeleteArticle(ctx context.Context, id int32) (*domain.DeleteKnowledgeBaseArticleResponse, error) {
	if err := s.repo.DeleteArticle(ctx, id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, &AppError{Status: 404, Message: "Статья не найдена"}
		}
		return nil, err
	}
	return &domain.DeleteKnowledgeBaseArticleResponse{Message: "Статья успешно удалена"}, nil
}
