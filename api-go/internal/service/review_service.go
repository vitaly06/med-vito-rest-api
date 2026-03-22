package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"med-vito/api-go/internal/repository"
)

var reviewAllowedRatings = []float64{1, 1.5, 2, 2.5, 3, 3.5, 4, 4.5, 5}

func isAllowedReviewRating(r float64) bool {
	for _, a := range reviewAllowedRatings {
		if math.Abs(r-a) < 1e-6 {
			return true
		}
	}
	return false
}

type ReviewService struct {
	repo *repository.ReviewPG
}

func NewReviewService(repo *repository.ReviewPG) *ReviewService {
	return &ReviewService{repo: repo}
}

func (s *ReviewService) SendReview(ctx context.Context, authorID int32, reviewedUserID int32, rating float64, text *string) (map[string]any, error) {
	if reviewedUserID <= 0 {
		return nil, &AppError{400, "Id оцениваемого продавца должен быть положительным числом"}
	}
	if !isAllowedReviewRating(rating) {
		return nil, &AppError{400, "Рейтинг должен быть одним из: 1, 1.5, 2, 2.5, 3, 3.5, 4, 4.5, 5"}
	}
	ok, err := s.repo.UserExists(ctx, reviewedUserID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, &AppError{400, "Продавец с таким id не найден"}
	}
	if authorID == reviewedUserID {
		return nil, &AppError{400, "Нельзя оставить отзыв самому себе"}
	}
	if err := s.repo.InsertReview(ctx, authorID, reviewedUserID, rating, text); err != nil {
		if repository.IsPgUniqueViolation(err) {
			return nil, &AppError{400, "Вы уже оставили отзыв этому пользователю"}
		}
		return nil, err
	}
	return map[string]any{"message": "Отзыв успешно оставлен"}, nil
}

func (s *ReviewService) GetUserReviews(ctx context.Context, userID int32) (map[string]any, error) {
	rows, err := s.repo.ListApprovedForUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	var sum float64
	for _, r := range rows {
		sum += r.Rating
	}
	var total float64
	if len(rows) > 0 {
		total = math.Round((sum/float64(len(rows)))*100) / 100
	}
	list := make([]map[string]any, 0, len(rows))
	for _, r := range rows {
		lt := r.CreatedAt.In(time.Local)
		list = append(list, map[string]any{
			"rating":   r.Rating,
			"text":     r.Text,
			"date":     fmt.Sprintf("%02d.%02d.%02d", lt.Day(), int(lt.Month()), lt.Year()%100),
			"fullName": r.FullName,
		})
	}
	return map[string]any{
		"totalRating":  total,
		"reviewsCount": len(rows),
		"reviews":      list,
	}, nil
}

func (s *ReviewService) ModerateReview(ctx context.Context, reviewID int32, status string) (map[string]any, error) {
	st := strings.TrimSpace(strings.ToUpper(status))
	if st != "APPROVED" && st != "DENIDED" {
		return nil, &AppError{400, "Неверный статус модерации. Доступные статуты: APPROVED, DENIDED"}
	}
	if err := s.repo.SetReviewModeration(ctx, reviewID, st); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, &AppError{404, "Отзыв для модерации не найден"}
		}
		return nil, err
	}
	if st == "APPROVED" {
		return map[string]any{"message": "Отзыв успешно опубликован"}, nil
	}
	return map[string]any{"message": "Отзыв успешно отклонён"}, nil
}

func (s *ReviewService) AllReviewsToModerate(ctx context.Context) ([]map[string]any, error) {
	rows, err := s.repo.ListModerateQueue(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]map[string]any, 0, len(rows))
	for _, r := range rows {
		out = append(out, map[string]any{
			"id":        r.ID,
			"rating":    r.Rating,
			"text":      r.Text,
			"createdAt": r.CreatedAt,
			"reviewedBy": map[string]any{
				"id":       r.ReviewedByID,
				"fullName": r.ReviewedByName,
			},
			"reviewedUser": map[string]any{
				"id":       r.ReviewedUserID,
				"fullName": r.ReviewedName,
			},
		})
	}
	return out, nil
}
