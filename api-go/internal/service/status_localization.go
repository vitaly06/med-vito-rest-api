package service

var dealStatusRu = map[string]string{
	"CREATED":   "Создана",
	"PAID":      "Оплачена",
	"SHIPPED":   "Отправлена",
	"DELIVERED": "Доставлена",
	"COMPLETED": "Завершена",
	"CANCELLED": "Отменена",
	"DISPUTE":   "Спор",
	"REFUNDED":  "Возврат",
}

var reservationStatusRu = map[string]string{
	"ACTIVE":              "Активен",
	"CANCELLED_BY_BUYER":  "Отменен покупателем",
	"CANCELLED_BY_SELLER": "Отменен продавцом",
	"EXPIRED":             "Истек",
	"DEAL_CREATED":        "Сделка создана",
	"COMPLETED":           "Завершен",
}

func localizeDealStatus(status string) string {
	if v, ok := dealStatusRu[status]; ok {
		return v
	}
	return status
}

func localizeReservationStatus(status string) string {
	if v, ok := reservationStatusRu[status]; ok {
		return v
	}
	return status
}
