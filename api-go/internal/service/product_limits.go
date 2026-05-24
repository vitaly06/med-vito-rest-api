package service

import (
	"strings"
)

const (
	freeMaxImages  = 8
	paidMaxImages  = 15
	freeDescMaxLen = 2000
	paidDescMaxLen = 7000
)

func validateProductContentLimits(imagesCount int, description string, paid bool) error {
	maxImages := freeMaxImages
	maxDesc := freeDescMaxLen
	if paid {
		maxImages = paidMaxImages
		maxDesc = paidDescMaxLen
	}

	if imagesCount > maxImages {
		if paid {
			return &AppError{400, "Для платного объявления можно загрузить не более 15 фото"}
		}
		return &AppError{400, "Для бесплатного объявления можно загрузить не более 8 фото"}
	}

	descLen := len([]rune(strings.TrimSpace(description)))
	if descLen > maxDesc {
		if paid {
			return &AppError{400, "Для платного объявления описание не должно превышать 7000 символов"}
		}
		return &AppError{400, "Для бесплатного объявления описание не должно превышать 2000 символов"}
	}
	return nil
}
