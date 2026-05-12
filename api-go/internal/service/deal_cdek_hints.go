package service

import (
	"fmt"
	"strings"

	"med-vito/api-go/internal/repository"
)

// buildCdekSellerHandoffNote — текст для продавца из тарифа и ПВЗ (без противоречия «дверь-дверь» vs «сдай в ПВЗ»).
func buildCdekSellerHandoffNote(deal repository.DealRow) string {
	tn := ""
	if deal.CDEKTariffName != nil {
		tn = strings.ToLower(strings.TrimSpace(*deal.CDEKTariffName))
	}
	toPvz := ""
	if deal.CDEKToPVZ != nil {
		toPvz = strings.TrimSpace(*deal.CDEKToPVZ)
	}
	fromPvz := ""
	if deal.CDEKFromPVZ != nil {
		fromPvz = strings.TrimSpace(*deal.CDEKFromPVZ)
	}

	hasDoor := strings.Contains(tn, "дверь")
	hasPvzWord := strings.Contains(tn, "пвз") || strings.Contains(tn, "склад") ||
		strings.Contains(tn, "постамат") || strings.Contains(tn, "пункт")

	if hasDoor && !hasPvzWord {
		if toPvz != "" {
			return fmt.Sprintf(
				"Тариф с доставкой до двери: при сдаче в CDEK согласуй адреса забора и вручения. В сделке также указан ПВЗ покупателя %s — если фактическая выдача там, проверь в ЛК CDEK, что тариф и точка согласованы.",
				toPvz,
			)
		}
		return "Тариф с доставкой до двери: при сдаче в CDEK согласуй адреса забора и вручения в личном кабинете или с менеджером."
	}

	if toPvz != "" {
		return fmt.Sprintf("Передай отправление в CDEK для выдачи в ПВЗ покупателя: %s.", toPvz)
	}
	if fromPvz != "" {
		return fmt.Sprintf("Сдай отправление в своём пункте CDEK: %s.", fromPvz)
	}
	return "Проверь в личном кабинете CDEK тариф и точки отправления/получения по этой сделке."
}

// buildCdekRegistrationHint — когда ждать UUID/трек (данные с нашей логики + CDEK API).
func buildCdekRegistrationHint(deal repository.DealRow) string {
	hasUUID := deal.CDEKOrderUUID != nil && strings.TrimSpace(*deal.CDEKOrderUUID) != ""
	track := ""
	if deal.CDEKTrackNumber != nil {
		track = strings.TrimSpace(*deal.CDEKTrackNumber)
	}

	if hasUUID && track == "" {
		return "Трек-номер присвоит CDEK после приёма отправления в сеть — обнови страницу позже или открой трекинг по ссылке, когда появится."
	}
	if hasUUID {
		return ""
	}

	switch deal.Status {
	case "CREATED":
		return "Заказ в CDEK создаётся после оплаты покупателем — здесь появятся UUID и трек из API CDEK."
	case "PAID":
		return "Ожидаем регистрацию в CDEK (автоматически). Обнови через 1–2 минуты. Если UUID пустой — проверь ключи CDEK на сервере и логи API."
	default:
		if deal.Status == "SHIPPED" || deal.Status == "DELIVERED" {
			return "UUID заказа CDEK в этой сделке не сохранён — расширенные данные из API CDEK могут быть недоступны."
		}
		return ""
	}
}
