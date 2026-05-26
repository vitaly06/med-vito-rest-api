package service

import (
	"fmt"
	"strings"

	"med-vito/api-go/internal/repository"
)

func buildCdekSellerHandoffNote(deal repository.DealRow) string {
	if deal.CDEKSellerHandoff == nil || strings.TrimSpace(*deal.CDEKSellerHandoff) == "" {
		return "После оплаты выбери способ передачи: сдать в ПВЗ СДЭК или вызвать курьера. Затем принеси посылку — сотрудник проверит и упакует."
	}
	switch strings.TrimSpace(*deal.CDEKSellerHandoff) {
	case "courier":
		addr := ""
		if deal.CDEKFromAddress != nil {
			addr = strings.TrimSpace(*deal.CDEKFromAddress)
		}
		if addr != "" {
			return fmt.Sprintf("Курьер СДЭК заберёт посылку по адресу: %s. Дождись приезда, передай груз после проверки.", addr)
		}
		return "Курьер СДЭК заберёт посылку с указанного адреса. После приёма в сеть появится трек-номер."
	case "pvz":
		if deal.CDEKFromPVZ != nil && strings.TrimSpace(*deal.CDEKFromPVZ) != "" {
			return fmt.Sprintf("Принеси посылку в пункт СДЭК %s. Сотрудник проверит содержимое, упакует и примет отправление.", strings.TrimSpace(*deal.CDEKFromPVZ))
		}
		return "Принеси посылку в выбранный пункт СДЭК — сотрудник проверит содержимое и упакует груз."
	default:
		return buildCdekSellerHandoffNoteLegacy(deal)
	}
}

func buildCdekSellerHandoffNoteLegacy(deal repository.DealRow) string {
	tn := ""
	if deal.CDEKTariffName != nil {
		tn = strings.ToLower(strings.TrimSpace(*deal.CDEKTariffName))
	}
	toPvz := ""
	if deal.CDEKToPVZ != nil {
		toPvz = strings.TrimSpace(*deal.CDEKToPVZ)
	}
	if strings.Contains(tn, "дверь") && toPvz != "" {
		return fmt.Sprintf("Тариф с доставкой до двери. ПВЗ покупателя в сделке: %s — уточни в ЛК СДЭК точку вручения.", toPvz)
	}
	if toPvz != "" {
		return fmt.Sprintf("Передай отправление в СДЭК для выдачи покупателю в ПВЗ: %s.", toPvz)
	}
	return "Проверь в личном кабинете СДЭК тариф и точки по этой сделке."
}

func buildCdekRegistrationHint(deal repository.DealRow) string {
	hasUUID := deal.CDEKOrderUUID != nil && strings.TrimSpace(*deal.CDEKOrderUUID) != ""
	track := ""
	if deal.CDEKTrackNumber != nil {
		track = strings.TrimSpace(*deal.CDEKTrackNumber)
	}

	if hasUUID && track == "" {
		return "Трек-номер присвоит СДЭК после приёма посылки в сеть (этап «Передача»). Обнови страницу позже."
	}
	if hasUUID {
		return ""
	}

	switch deal.Status {
	case "CREATED":
		return "После оплаты продавец оформит заявку в СДЭК (получатель и вес уже указаны при создании сделки)."
	case "PAID":
		if deal.CDEKSellerHandoff == nil || strings.TrimSpace(*deal.CDEKSellerHandoff) == "" {
			return "Продавец должен выбрать способ передачи (ПВЗ или курьер) — затем создастся заказ в СДЭК."
		}
		return "Заказ в СДЭК создаётся после выбора способа передачи. Если UUID пустой — обнови через минуту или проверь ключи CDEK."
	default:
		if deal.Status == "SHIPPED" || deal.Status == "DELIVERED" {
			return "UUID заказа СДЭК не сохранён — трекинг может быть ограничен."
		}
		return ""
	}
}
