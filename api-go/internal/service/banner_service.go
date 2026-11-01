package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"med-vito/api-go/internal/config"
	mailpkg "med-vito/api-go/internal/pkg/mail"
	"med-vito/api-go/internal/pkg/s3client"
	"med-vito/api-go/internal/repository"
)

var validBannerPlaces = map[string]bool{
	"PRODUCT_FEED": true, "PROFILE": true, "FAVORITES": true, "CHATS": true,
}

var validBannerModerate = map[string]bool{
	"MODERATE": true, "APPROVED": true, "DENIDED": true,
}

type BannerService struct {
	repo    *repository.BannerPG
	s3      *s3client.Client
	support *SupportService
	cfg     config.Config
}

func NewBannerService(repo *repository.BannerPG, s3 *s3client.Client, support *SupportService, cfg config.Config) *BannerService {
	return &BannerService{repo: repo, s3: s3, support: support, cfg: cfg}
}

func (s *BannerService) sendMail(to, subject, html string) {
	if strings.TrimSpace(s.cfg.SMTPHost) == "" || strings.TrimSpace(to) == "" {
		return
	}
	_ = mailpkg.SendHTMLSmart(s.cfg.SMTPHost, s.cfg.SMTPPort, s.cfg.SMTPUser, s.cfg.SMTPPassword, s.cfg.SMTPFrom, to, subject, html, s.cfg.SMTPSecure, s.cfg.SMTPTLSInsecure)
}

func (s *BannerService) Create(ctx context.Context, userID int32, isAdmin bool, file UploadedFile, place, navigateToURL, name string) (*repository.BannerRow, error) {
	if s.s3 == nil {
		return nil, &AppError{500, "S3 не настроен"}
	}
	place = strings.TrimSpace(place)
	if !validBannerPlaces[place] {
		return nil, &AppError{400, "Некорректное значение места размещения"}
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, &AppError{400, "Текст баннера обязателен для заполнения"}
	}
	navigateToURL = strings.TrimSpace(navigateToURL)
	if navigateToURL == "" {
		return nil, &AppError{400, "Navigate url обязателен для заполнения"}
	}
	if len(file.Body) == 0 {
		return nil, &AppError{400, "Изображение баннера обязательно"}
	}
	state := "MODERATE"
	if isAdmin {
		state = "APPROVED"
	}
	url, err := s.s3.Upload(ctx, "banners", file.Name, file.ContentType, file.Body)
	if err != nil {
		return nil, &AppError{400, "Ошибка загрузки в S3: " + err.Error()}
	}
	return s.repo.Insert(ctx, name, url, place, state, navigateToURL, userID)
}

func (s *BannerService) FindRandom(ctx context.Context, limit int) ([]repository.BannerRow, error) {
	return s.repo.FindRandomApproved(ctx, limit)
}

func (s *BannerService) FindAll(ctx context.Context, place *string) ([]repository.BannerRow, error) {
	if place != nil && *place != "" {
		if !validBannerPlaces[*place] {
			return nil, &AppError{400, "Некорректное значение места размещения"}
		}
		return s.repo.FindByPlaceApproved(ctx, *place)
	}
	return s.repo.FindAllApproved(ctx)
}

func (s *BannerService) FindOne(ctx context.Context, id int32) (*repository.BannerRow, error) {
	b, err := s.repo.GetByID(ctx, id)
	if errors.Is(err, repository.ErrBannerNotFound) {
		return nil, &AppError{404, "Баннер с ID " + strconv.Itoa(int(id)) + " не найден"}
	}
	return b, err
}

func (s *BannerService) Update(ctx context.Context, id int32, requesterID int32, isAdmin bool, file *UploadedFile, place, navigateToURL, name *string) (*repository.BannerRow, error) {
	if s.s3 == nil {
		return nil, &AppError{500, "S3 не настроен"}
	}
	existing, err := s.repo.GetByID(ctx, id)
	if errors.Is(err, repository.ErrBannerNotFound) {
		return nil, &AppError{404, "Баннер с ID " + strconv.Itoa(int(id)) + " не найден"}
	}
	if err != nil {
		return nil, err
	}
	if !isAdmin && existing.UserID != requesterID {
		return nil, &AppError{403, "Нет прав на редактирование баннера"}
	}
	if place != nil && *place != "" && !validBannerPlaces[*place] {
		return nil, &AppError{400, "Некорректное значение места размещения"}
	}
	var photoURL *string
	if file != nil && len(file.Body) > 0 {
		if err := s.s3.DeleteByURL(ctx, existing.PhotoURL); err != nil {
			return nil, err
		}
		u, err := s.s3.Upload(ctx, "banners", file.Name, file.ContentType, file.Body)
		if err != nil {
			return nil, &AppError{400, "Ошибка загрузки в S3: " + err.Error()}
		}
		photoURL = &u
	}
	var pl, nav, nm *string
	if place != nil && *place != "" {
		pl = place
	}
	if navigateToURL != nil {
		nav = navigateToURL
	}
	if name != nil && strings.TrimSpace(*name) != "" {
		x := strings.TrimSpace(*name)
		nm = &x
	}
	var moderateState *string
	if !isAdmin {
		mod := "MODERATE"
		moderateState = &mod
	}
	return s.repo.Update(ctx, id, photoURL, pl, nav, nm, moderateState)
}

func (s *BannerService) Toggle(ctx context.Context, id int32, requesterID int32, isAdmin bool) (*repository.BannerRow, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if errors.Is(err, repository.ErrBannerNotFound) {
		return nil, &AppError{404, "Баннер с ID " + strconv.Itoa(int(id)) + " не найден"}
	}
	if err != nil {
		return nil, err
	}
	if !isAdmin && existing.UserID != requesterID {
		return nil, &AppError{403, "Нет прав на изменение баннера"}
	}
	newStatus := "DENIDED"
	if existing.ModerateState == "DENIDED" {
		if isAdmin {
			newStatus = "APPROVED"
		} else {
			newStatus = "MODERATE"
		}
	} else if existing.ModerateState == "APPROVED" || existing.ModerateState == "MODERATE" {
		newStatus = "DENIDED"
	}
	return s.repo.Update(ctx, id, nil, nil, nil, nil, &newStatus)
}

func (s *BannerService) Remove(ctx context.Context, id int32, requesterID int32, isAdmin bool) (map[string]string, error) {
	if s.s3 == nil {
		return nil, &AppError{500, "S3 не настроен"}
	}
	b, err := s.repo.GetByID(ctx, id)
	if errors.Is(err, repository.ErrBannerNotFound) {
		return nil, &AppError{404, "Баннер с ID " + strconv.Itoa(int(id)) + " не найден"}
	}
	if err != nil {
		return nil, err
	}
	if !isAdmin && b.UserID != requesterID {
		return nil, &AppError{403, "Нет прав на удаление баннера"}
	}
	_ = s.s3.DeleteByURL(ctx, b.PhotoURL)
	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrBannerNotFound) {
			return nil, &AppError{404, "Баннер с ID " + strconv.Itoa(int(id)) + " не найден"}
		}
		return nil, err
	}
	return map[string]string{"message": "Баннер успешно удалён"}, nil
}

func (s *BannerService) RegisterView(ctx context.Context, bannerID int32, userID *int32, ip string) (*repository.BannerViewRow, error) {
	if _, err := s.repo.GetByID(ctx, bannerID); errors.Is(err, repository.ErrBannerNotFound) {
		return nil, &AppError{404, "Баннер с ID " + strconv.Itoa(int(bannerID)) + " не найден"}
	} else if err != nil {
		return nil, err
	}
	since := time.Now().Add(-24 * time.Hour)
	existing, err := s.repo.FindRecentView(ctx, bannerID, since, userID, ip)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return existing, nil
	}
	var ipPtr *string
	if strings.TrimSpace(ip) != "" {
		x := strings.TrimSpace(ip)
		ipPtr = &x
	}
	return s.repo.InsertView(ctx, bannerID, userID, ipPtr)
}

func (s *BannerService) GetBannerStats(ctx context.Context, bannerID int32) (map[string]any, error) {
	b, err := s.repo.GetByID(ctx, bannerID)
	if errors.Is(err, repository.ErrBannerNotFound) {
		return nil, &AppError{404, "Баннер с ID " + strconv.Itoa(int(bannerID)) + " не найден"}
	}
	if err != nil {
		return nil, err
	}
	total, err := s.repo.CountViews(ctx, bannerID)
	if err != nil {
		return nil, err
	}
	uniq, err := s.repo.CountDistinctViewUsers(ctx, bannerID)
	if err != nil {
		return nil, err
	}
	since := time.Now().AddDate(0, 0, -30)
	byDay, err := s.repo.ViewsByDay(ctx, bannerID, since)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"banner":       b,
		"totalViews":   total,
		"uniqueUsers":  uniq,
		"viewsByDay":   byDay,
	}, nil
}

func (s *BannerService) GetUserBannersStats(ctx context.Context, userID int32) ([]repository.UserBannerStat, error) {
	return s.repo.UserBannersStats(ctx, userID)
}

func (s *BannerService) ModerateBanner(ctx context.Context, id int32, status string) (map[string]string, error) {
	status = strings.TrimSpace(status)
	if !validBannerModerate[status] {
		return nil, &AppError{400, "Неверный статус модерации. Доступные статусы: MODERATE, APPROVED, DENIDED"}
	}
	existing, err := s.repo.GetByID(ctx, id)
	if errors.Is(err, repository.ErrBannerNotFound) {
		return nil, &AppError{404, "Баннер для модерации не найден"}
	} else if err != nil {
		return nil, err
	}

	switch status {
	case "APPROVED":
		var price float64
		switch existing.Place {
		case "PRODUCT_FEED":
			price = 500
		case "PROFILE":
			price = 800
		case "CHATS":
			price = 1200
		case "FAVORITES":
			price = 1000
		default:
			price = 500
		}

		res, err := s.repo.ApproveAndDeduct(ctx, id, price)
		if err != nil {
			var ie *repository.BannerInsufficientError
			if errors.As(err, &ie) {
				return nil, &AppError{400, fmt.Sprintf("У пользователя недостаточно средств (требуется %.0f ₽, доступно %.0f ₽)", ie.Required, ie.Available)}
			}
			return nil, err
		}

		if !res.AlreadyApproved {
			if s.support != nil && res.UserID > 0 {
				chatMsg := fmt.Sprintf("✅ Ваш баннер «%s» успешно одобрен и опубликован.\nС баланса списано: %.0f ₽.", res.BannerName, res.Price)
				go s.support.NotifyUserBilling(context.Background(), res.UserID, chatMsg)
			}
			if res.UserEmail != "" {
				subject := fmt.Sprintf("Баннер «%s» одобрен и опубликован — TorguiSam.ru", res.BannerName)
				body := fmt.Sprintf("<p>Здравствуйте!</p><p>Ваш рекламный баннер <b>«%s»</b> успешно прошёл модерацию и опубликован.</p><p>С вашего баланса списано: <b>%.0f ₽</b>.</p>", res.BannerName, res.Price)
				go s.sendMail(res.UserEmail, subject, body)
			}
		}
		return map[string]string{"message": "Баннер успешно опубликован"}, nil

	case "DENIDED":
		if err := s.repo.SetModerateState(ctx, id, "DENIDED"); err != nil {
			return nil, err
		}
		bannerName, userID, userEmail, _ := s.repo.GetBannerOwnerInfo(ctx, id)
		if s.support != nil && userID > 0 {
			chatMsg := fmt.Sprintf("❌ Ваш баннер «%s» был отклонён модератором.", bannerName)
			go s.support.NotifyUserBilling(context.Background(), userID, chatMsg)
		}
		if userEmail != "" {
			subject := fmt.Sprintf("Баннер «%s» отклонён — TorguiSam.ru", bannerName)
			body := fmt.Sprintf("<p>Здравствуйте!</p><p>Ваш рекламный баннер <b>«%s»</b> был отклонён модератором.</p>", bannerName)
			go s.sendMail(userEmail, subject, body)
		}
		return map[string]string{"message": "Баннер успешно отклонён"}, nil

	default:
		if err := s.repo.SetModerateState(ctx, id, "MODERATE"); err != nil {
			return nil, err
		}
		return map[string]string{"message": "Баннер возвращён на модерацию"}, nil
	}
}

func (s *BannerService) AllBannersToModerate(ctx context.Context) ([]repository.BannerModerateRow, error) {
	return s.repo.ListModerateQueue(ctx)
}
