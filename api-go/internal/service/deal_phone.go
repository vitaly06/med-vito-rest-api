package service

import (
	"fmt"
	"strings"
	"unicode"
)

// PhoneForCdekAPI — нормализует номер в формат +7..., который ждёт CDEK (после проверки на «подставной»).
func PhoneForCdekAPI(phone string) (string, error) {
	if isDealPhoneSynthetic(phone) {
		return "", fmt.Errorf("подставной или пустой телефон")
	}
	var b strings.Builder
	for _, r := range phone {
		if unicode.IsDigit(r) {
			b.WriteRune(r)
		}
	}
	ds := b.String()
	switch len(ds) {
	case 10:
		if ds[0] == '9' {
			return "+7" + ds, nil
		}
	case 11:
		if ds[0] == '8' {
			return "+7" + ds[1:], nil
		}
		if ds[0] == '7' {
			return "+" + ds, nil
		}
	}
	return "", fmt.Errorf("не удалось разобрать телефон для CDEK")
}

// isDealPhoneSynthetic — true, если в "phoneNumber" нет нормального РФ-мобильного (часто VK_/oauth после OAuth).
func isDealPhoneSynthetic(phone string) bool {
	p := strings.TrimSpace(phone)
	if p == "" {
		return true
	}
	up := strings.ToUpper(p)
	for _, pref := range []string{"VK_", "MAX_", "OAUTH_", "TELEGRAM_", "GOOGLE_", "YANDEX_"} {
		if strings.HasPrefix(up, pref) {
			return true
		}
	}
	var b strings.Builder
	for _, r := range p {
		if unicode.IsDigit(r) {
			b.WriteRune(r)
		}
	}
	ds := b.String()
	switch len(ds) {
	case 11:
		if ds[0] == '7' || ds[0] == '8' {
			return false
		}
	case 10:
		if ds[0] == '9' {
			return false
		}
	}
	return true
}
