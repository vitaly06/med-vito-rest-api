package service

import (
	"strings"

	"med-vito/api-go/internal/repository"
)

// Этапы доставки СДЭК по сценарию: оформление → передача → логистика → оповещение → вручение.
func buildCdekDeliveryStages(deal repository.DealRow) []map[string]any {
	hasRoute := dealHasCdekAutoRoute(&deal)
	if !hasRoute {
		return nil
	}

	paid := deal.Status != "CREATED" && deal.Status != "CANCELLED"
	hasOrder := deal.CDEKOrderUUID != nil && strings.TrimSpace(*deal.CDEKOrderUUID) != ""
	hasHandoff := deal.CDEKSellerHandoff != nil && strings.TrimSpace(*deal.CDEKSellerHandoff) != ""
	shipped := deal.Status == "SHIPPED" || deal.Status == "DELIVERED" || deal.Status == "COMPLETED"
	delivered := deal.Status == "DELIVERED" || deal.Status == "COMPLETED"

	track := ""
	if deal.CDEKTrackNumber != nil {
		track = strings.TrimSpace(*deal.CDEKTrackNumber)
	}
	hasTrack := track != ""

	cdekStatus := ""
	if deal.CDEKStatus != nil {
		cdekStatus = strings.ToUpper(strings.TrimSpace(*deal.CDEKStatus))
	}
	atPickup := cdekStatusIndicatesReadyForPickup(cdekStatus)
	inTransit := cdekStatusIndicatesInTransit(cdekStatus) || (hasTrack && shipped && !atPickup)

	stages := []struct {
		key, title, desc string
		done, active     bool
	}{
		{
			key:   "order",
			title: "Оформление",
			desc:  "Заявка на отправление: получатель, вес и тариф. После оплаты продавец оформит передачу в СДЭК.",
			done:  paid,
			active: deal.Status == "CREATED" || (deal.Status == "PAID" && !hasHandoff),
		},
		{
			key:   "handoff",
			title: "Передача",
			desc:  handoffStageDescription(deal),
			done:  hasHandoff && (hasOrder || shipped),
			active: deal.Status == "PAID" && paid && !shipped,
		},
		{
			key:   "logistics",
			title: "Логистика",
			desc:  "Посылка в пути: сортировочный узел, доставка по выбранному тарифу (авиа/авто/жд).",
			done:  inTransit || atPickup || delivered,
			active: shipped && !delivered && !atPickup,
		},
		{
			key:   "notification",
			title: "Оповещение",
			desc:  notificationStageDescription(),
			done:  atPickup || delivered,
			active: hasTrack && shipped && !delivered && !atPickup,
		},
		{
			key:   "delivery",
			title: "Вручение",
			desc:  deliveryStageDescription(),
			done:  delivered,
			active: atPickup || (shipped && !delivered),
		},
	}

	out := make([]map[string]any, 0, len(stages))
	for _, st := range stages {
		status := "pending"
		if st.done {
			status = "done"
		} else if st.active {
			status = "active"
		}
		out = append(out, map[string]any{
			"key":         st.key,
			"title":       st.title,
			"description": st.desc,
			"status":      status,
		})
	}
	return out
}

func handoffStageDescription(deal repository.DealRow) string {
	if deal.CDEKSellerHandoff != nil {
		switch strings.TrimSpace(*deal.CDEKSellerHandoff) {
		case "courier":
			return "Продавец вызвал курьера СДЭК: курьер заберёт посылку по указанному адресу, проверит и упакует."
		case "pvz":
			if deal.CDEKFromPVZ != nil && strings.TrimSpace(*deal.CDEKFromPVZ) != "" {
				return "Продавец сдаёт посылку в пункт СДЭК " + strings.TrimSpace(*deal.CDEKFromPVZ) + " — сотрудник примет и отправит в сеть."
			}
			return "Продавец приносит посылку в пункт выдачи СДЭК — сотрудник проверяет содержимое и упаковывает."
		}
	}
	return "Продавец сдаёт посылку в ПВЗ СДЭК или вызывает курьера на адрес — после приёма в сеть появится трек-номер."
}

func notificationStageDescription() string {
	return "Когда груз прибудет в пункт выдачи, СДЭК пришлёт получателю SMS или push с кодом для получения."
}

func deliveryStageDescription() string {
	return "Получатель забирает посылку в ПВЗ по паспорту или коду из SMS. Затем подтверждает получение в сделке."
}

func cdekStatusIndicatesInTransit(status string) bool {
	if status == "" {
		return false
	}
	for _, part := range []string{
		"TRANSIT", "DELIVERING", "WAREHOUSE", "ACCEPTED", "RECEIVED", "SENT",
		"IN_TRANSIT", "ON_WAY", "COURIER",
	} {
		if strings.Contains(status, part) {
			return true
		}
	}
	return false
}

func cdekStatusIndicatesReadyForPickup(status string) bool {
	if status == "" {
		return false
	}
	for _, part := range []string{
		"READY_FOR_PICKUP", "READY_TO_PICKUP", "POSTOMAT", "PVZ", "PICKUP",
		"AWAITING", "STORAGE",
	} {
		if strings.Contains(status, part) {
			return true
		}
	}
	return strings.Contains(status, "DELIVERED") && !strings.Contains(status, "RECIPIENT")
}
